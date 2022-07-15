package user

import (
	"context"
	"fmt"
	"github.com/devpies/saas-core/internal/user/clients"
	"github.com/nats-io/nats.go"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/devpies/saas-core/internal/user/config"
	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/handler"
	"github.com/devpies/saas-core/internal/user/repository"
	"github.com/devpies/saas-core/internal/user/res"
	"github.com/devpies/saas-core/internal/user/service"
	"github.com/devpies/saas-core/pkg/log"
	"github.com/devpies/saas-core/pkg/msg"

	"go.uber.org/zap"
)

// Run contains the app setup.
func Run() error {
	var (
		cfg     config.Config
		logger  *zap.Logger
		logPath = "log/out.log"
		err     error
	)

	cfg, err = config.NewConfig()
	if err != nil {
		return err
	}

	if cfg.Web.Production {
		logger, err = log.NewProductionLogger(logPath)
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	defer logger.Sync()

	pg, Close, err := db.NewPostgresDatabase(logger, cfg)
	if err != nil {
		return err
	}
	defer Close()

	// Execute latest migration in production.
	if cfg.Web.Production {
		if err = res.MigrateUp(pg.URL.String()); err != nil {
			logger.Error("error connecting to user database", zap.Error(err))
			return err
		}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	jetstream := msg.NewStreamContext(logger, shutdown, cfg.Nats.Address, cfg.Nats.Port)

	_ = jetstream.Create(msg.StreamMemberships)

	ctx := context.Background()

	dynamoDBClient := clients.NewDynamoDBClient(ctx, cfg.Cognito.Region)
	cognitoClient := clients.NewCognitoClient(ctx, cfg.Cognito.Region)

	// Initialize 3-layered architecture.
	inviteRepo := repository.NewInviteRepository(logger, pg)
	userRepo := repository.NewUserRepository(logger, pg)
	seatRepo := repository.NewSeatRepository(logger, pg)
	connections := repository.NewConnectionRepository(dynamoDBClient, cfg.Dynamodb.ConnectionTable)

	userService := service.NewUserService(logger, userRepo, seatRepo, cognitoClient, connections, cfg.Cognito.SharedUserPoolID)
	inviteService := service.NewInviteService(logger, inviteRepo)

	userHandler := handler.NewUserHandler(logger, userService)
	inviteHandler := handler.NewInviteHandler(logger, inviteService)
	opts := []nats.SubOpt{nats.DeliverAll(), nats.ManualAck()}

	go func() {
		jetstream.Listen(
			string(msg.TypeTenantIdentityCreated),
			msg.SubjectTenantIdentityCreated,
			"tenant_created_consumer",
			userService.AddAdminUserFromEvent,
			opts...)
	}()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, cfg.Cognito.Region, cfg.Cognito.SharedUserPoolID, userHandler, inviteHandler),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting user service on %s:%s", cfg.Web.Address, cfg.Web.Port))
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err = <-serverErrors:
		logger.Error("error on startup", zap.Error(err))
		return err
	case sig := <-shutdown:
		logger.Info(fmt.Sprintf("Start shutdown due to %s signal", sig))

		// Give on going tasks a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			err = srv.Close()
		}

		switch {
		case sig == syscall.SIGSTOP:
			logger.Error("error on integrity issue caused shutdown", zap.Error(err))
			return err
		case err != nil:
			logger.Error("error on gracefully shutdown", zap.Error(err))
			return err
		}
	}

	return err
}

package tenant

import (
	"context"
	"fmt"
	"github.com/devpies/saas-core/internal/tenant/clients"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/devpies/saas-core/internal/tenant/config"
	"github.com/devpies/saas-core/internal/tenant/handler"
	"github.com/devpies/saas-core/internal/tenant/repository"
	"github.com/devpies/saas-core/internal/tenant/service"
	"github.com/devpies/saas-core/pkg/log"
	"github.com/devpies/saas-core/pkg/msg"

	"github.com/nats-io/nats.go"
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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	// Initialize NATS JetStream.
	js := msg.NewStreamContext(logger, shutdown, cfg.Nats.Address, cfg.Nats.Port)
	opts := []nats.SubOpt{nats.DeliverAll(), nats.ManualAck()}

	ctx := context.Background()

	dynamoDBClient := clients.NewDynamoDBClient(ctx, cfg.Cognito.Region)
	cognitoClient := clients.NewCognitoClient(ctx, cfg.Cognito.Region)

	// Initialize 3-layered architecture.
	tenantRepository := repository.NewTenantRepository(dynamoDBClient, cfg.Dynamodb.TenantTable)
	siloConfigRepository := repository.NewSiloConfigRepository(dynamoDBClient, cfg.Dynamodb.ConfigTable)
	authInfoRepository := repository.NewAuthInfoRepository(logger, dynamoDBClient, cfg.Dynamodb.AuthTable)
	connectionRepository := repository.NewConnectionRepository(dynamoDBClient, cfg.Dynamodb.ConnectionTable)
	tenantService := service.NewTenantService(logger, js, cfg.Cognito.SharedUserPoolID, cognitoClient, tenantRepository, connectionRepository)
	authInfoService := service.NewAuthInfoService(logger, authInfoRepository, cfg.Cognito.Region)
	siloConfigService := service.NewSiloConfigService(logger, siloConfigRepository)

	tenantHandler := handler.NewTenantHandler(logger, tenantService)
	authInfoHandler := handler.NewAuthInfoHandler(logger, authInfoService)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("listener panic: %v", r)
				logger.Error(fmt.Sprintf("%s", debug.Stack()), zap.Error(err))
			}
		}()

		js.Listen(
			string(msg.TypeTenantRegistered),
			msg.SubjectTenantRegistered,
			"tenant_consumer",
			tenantService.CreateTenantFromEvent,
			opts...,
		)

		js.Listen(
			string(msg.TypeTenantSiloed),
			msg.SubjectTenantSiloed,
			"tenant_silo_consumer",
			siloConfigService.StoreConfigFromEvent,
			opts...,
		)
	}()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, cfg.Cognito.Region, cfg.Cognito.UserPoolID, tenantHandler, authInfoHandler),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting tenant service on %s:%s", cfg.Web.Address, cfg.Web.Port))
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

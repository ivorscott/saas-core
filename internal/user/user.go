package user

import (
	"context"
	"fmt"
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

	"github.com/ardanlabs/conf"
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

	if err = conf.Parse(os.Args[1:], "USER", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("USER", &cfg)
			if err != nil {
				logger.Error("error generating config usage", zap.Error(err))
				return err
			}
			fmt.Println(usage)
			return nil
		}
		logger.Error("error parsing config", zap.Error(err))
		return err
	}

	if cfg.Web.Production {
		logger, err = log.NewProductionLogger(logPath)
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		logger.Error("error creating logger", zap.Error(err))
		return err
	}
	defer logger.Sync()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	pg, Close, err := db.NewPostgresDatabase(logger, cfg)
	if err != nil {
		return err
	}
	defer Close()

	// Execute latest migration.
	if err = res.MigrateUp(pg.URL.String()); err != nil {
		logger.Error("error connecting to admin database", zap.Error(err))
		return err
	}

	jetStream := msg.NewStreamContext(logger, shutdown, cfg.Nats.Address, cfg.Nats.Port)

	_ = jetStream.Create(msg.StreamMemberships)

	// Initialize 3-layered architecture.
	inviteRepo := repository.NewInviteRepository(logger, pg)
	userRepo := repository.NewUserRepository(logger, pg)
	teamRepo := repository.NewTeamRepository(logger, pg)
	membershipRepo := repository.NewMembershipRepository(logger, pg)
	projectRepo := repository.NewProjectRepository(logger, pg)

	userService := service.NewUserService(logger, userRepo)
	teamService := service.NewTeamService(logger, teamRepo, inviteRepo)
	membershipService := service.NewMembershipService(logger, membershipRepo)
	projectService := service.NewProjectService(logger, projectRepo)
	inviteService := service.NewInviteService(logger, inviteRepo)

	userHandler := handler.NewUserHandler(logger, userService)
	teamHandler := handler.NewTeamHandler(logger, jetStream, cfg.Sendgrid.APIKey, teamService, projectService, membershipService, inviteService)
	membershipHandler := handler.NewMembershipHandler(logger, membershipService)

	go func() {
		// Listen to project events to save a redundant copy in the database.
	}()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, userHandler, teamHandler, membershipHandler, cfg),
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

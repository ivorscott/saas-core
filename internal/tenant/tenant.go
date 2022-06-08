package tenant

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/devpies/saas-core/internal/tenant/config"
	"github.com/devpies/saas-core/internal/tenant/db"
	"github.com/devpies/saas-core/internal/tenant/handler"
	"github.com/devpies/saas-core/internal/tenant/repository"
	"github.com/devpies/saas-core/internal/tenant/service"
	"github.com/devpies/saas-core/pkg/log"

	"github.com/ardanlabs/conf"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

// Run contains the app setup.
func Run() error {
	var (
		cfg      config.Config
		logger   *zap.Logger
		dbClient *dynamodb.Client
		logPath  = "log/out.log"
		err      error
	)

	if err = conf.Parse(os.Args[1:], "TENANT", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("TENANT", &cfg)
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

	ctx := context.Background()

	if cfg.Web.Production {
		logger, err = log.NewProductionLogger(logPath)
		dbClient = db.NewProductionDynamoDBClient(ctx)
	} else {
		logger, err = zap.NewDevelopment()
		dbClient = db.NewDevelopmentDynamoDBClient(ctx, cfg.Dynamodb.Port)
	}
	if err != nil {
		logger.Error("error creating logger", zap.Error(err))
		return err
	}
	defer logger.Sync()

	// Initialize 3-layered architecture.
	tenantRepository := repository.NewTenantRepository(dbClient)
	tenantConfigRepository := repository.NewTenantConfigRepository(dbClient)
	authInfoRepository := repository.NewAuthInfoRepository(dbClient)

	tenantService := service.NewTenantService(logger, cfg, tenantRepository, tenantConfigRepository, authInfoRepository)
	tenantHandler := handler.NewTenantHandler(logger, tenantService)

	go func() {
		// listen for messages
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, tenantHandler, cfg),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting tenant api on %s:%s", cfg.Web.Address, cfg.Web.Port))
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
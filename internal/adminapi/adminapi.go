package adminapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devpies/core/internal/adminapi/config"
	"github.com/devpies/core/internal/adminapi/db"
	"github.com/devpies/core/internal/adminapi/handler"
	"github.com/devpies/core/internal/adminapi/res"
	"github.com/devpies/core/internal/adminapi/service"
	"github.com/devpies/core/pkg/log"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/ardanlabs/conf"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"go.uber.org/zap"
)

var logPath = "log/out.log"

func Run() error {
	var (
		cfg    config.Config
		logger *zap.Logger
		err    error
	)

	// Initialize configuration.
	if err = conf.Parse(os.Args[1:], "ADMIN", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("ADMIN", &cfg)
			if err != nil {
				return fmt.Errorf("error generating config usage: %w", err)
			}
			fmt.Println(usage)
			return nil
		}
		return fmt.Errorf("error parsing config: %w", err)
	}

	// Initialize logger.
	if cfg.Web.Production {
		logger, err = log.NewProductionLogger(logPath)
		if err != nil {
			return err
		}
	} else {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	}
	defer logger.Sync()

	// Initialize AWS clients.
	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	cognitoClient := cip.NewFromConfig(awsCfg)

	// Initialize admin database.
	repo, err := db.NewPostgresRepository(cfg)
	if err != nil {
		return err
	}
	defer repo.Close()

	// Execute latest migration.
	if err = res.MigrateUp(repo.URL.String()); err != nil {
		logger.Fatal("", zap.Error(err))
	}

	// Initialize a session manager..
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = postgresstore.New(repo.DB.DB)

	authService := service.NewAuthService(logger, cfg, cognitoClient)
	authHandler := handler.NewAuth(logger, authService, session)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.BackendPort),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      API(shutdown, logger, authHandler),
	}
	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Starting server...")
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error on startup : %w", err)
	case sig := <-shutdown:
		logger.Info(fmt.Sprintf("Start shutdown due to %s signal", sig))

		// give on going tasks a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			err = srv.Close()
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.New("could not stop server gracefully")
		}
	}

	return err
}

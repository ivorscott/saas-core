package adminclient

import (
	"context"
	"embed"
	"errors"
	"fmt"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/devpies/core/internal/adminclient/handler"
	"github.com/devpies/core/internal/adminclient/res"
	"github.com/devpies/core/internal/adminclient/service"
	"os/signal"
	"syscall"

	"github.com/devpies/core/internal/adminclient/db"

	"github.com/devpies/core/internal/adminclient/config"
	"github.com/devpies/core/internal/adminclient/render"
	"github.com/devpies/core/internal/adminclient/webpage"
	"github.com/devpies/core/pkg/log"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/ardanlabs/conf"
	"go.uber.org/zap"

	"io/fs"
	"net/http"
	"os"
	"time"
)

func Run(staticFS embed.FS) error {
	var (
		cfg     config.Config
		logger  *zap.Logger
		logPath = "log/out.log"
		session *scs.SessionManager
		err     error
	)

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

	// Initialize AWS clients.
	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	cognitoClient := cip.NewFromConfig(awsCfg)

	// Initialize 3-layered architecture.
	authService := service.NewAuthService(logger, cfg, cognitoClient)
	authHandler := handler.NewAuthHandler(logger, authService)

	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = postgresstore.New(repo.DB.DB)

	templateFS, err := fs.Sub(staticFS, "static/templates")
	if err != nil {
		logger.Error("", zap.Error(err))
	}
	assets, err := fs.Sub(staticFS, "static/assets")
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	renderEngine := render.New(logger, cfg, templateFS)
	pages := webpage.New(logger, cfg, renderEngine, session)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.FrontendPort),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      API(logger, shutdown, assets, pages, authHandler),
	}

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

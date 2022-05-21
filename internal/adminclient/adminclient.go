package adminclient

import (
	"embed"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/ardanlabs/conf"
	"github.com/devpies/core/internal/adminclient/config"
	"github.com/devpies/core/internal/adminclient/render"
	"github.com/devpies/core/internal/adminclient/webpage"
	"github.com/devpies/core/pkg/log"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
	"os"
	"time"
)

var cfg config.Config
var logPath = "log/out.log"
var session *scs.SessionManager

func Run(staticFS embed.FS) error {
	var (
		logger *zap.Logger
		err    error
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

	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	templateFS, err := fs.Sub(staticFS, "static/templates")
	if err != nil {
		logger.Error("", zap.Error(err))
	}
	assets, err := fs.Sub(staticFS, "static/assets")
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	renderEngine := render.New(logger, cfg, templateFS)
	pages := webpage.New(logger, cfg, renderEngine)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.FrontendPort),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      API(assets, pages),
	}
	serverErrors := make(chan error, 1)

	logger.Info("Starting server...")
	serverErrors <- srv.ListenAndServe()

	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error on startup : %w", err)
	default:
	}

	return err
}

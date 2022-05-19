package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ardanlabs/conf"

	"github.com/alexedwards/scs/v2"
	"github.com/devpies/core/internal/admin"
	"github.com/devpies/core/internal/admin/webapp"
	"github.com/devpies/core/internal/admin/webapp/render"
	"github.com/devpies/core/pkg/log"
)

//go:embed static
var staticFS embed.FS
var cfg admin.Config
var logPath = "log/out.log"
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	var err error

	logger, Sync := log.NewLoggerOrPanic(logPath)
	defer Sync()

	if err := conf.Parse(os.Args[1:], "ADMIN", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("ADMIN", &cfg)
			if err != nil {
				return fmt.Errorf("error generating config usage: %w", err)
			}
			fmt.Println(usage)
			return nil
		}
		return fmt.Errorf("error parsing config: %w", err)
	}

	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	templateFS, err := fs.Sub(staticFS, "static/templates")

	if err != nil {
		logger.Error("", zap.Error(err))
	}
	renderEngine := render.New(logger, cfg, templateFS)
	app := webapp.New(logger, cfg, renderEngine)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.FrontendPort),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      API(shutdown, logger, staticFS, app),
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

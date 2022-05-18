package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/devpies/core/admin/pkg/config"
	"github.com/devpies/core/admin/pkg/render"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ardanlabs/conf"

	"github.com/devpies/core/pkg/log"
)

//go:embed static/templates
var templateFS embed.FS
var cfg config.Config
var logPath = ".log/out.log"
var session *scs.SessionManager

type templateCache map[string]*template.Template

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

	if err := conf.Parse(os.Args[1:], "API", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("API", &cfg)
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

	templateCache := make(map[string]*template.Template)
	templateEngine := render.New(logger, cfg, templateCache, templateFS)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      API(shutdown, logger, templateEngine, templateFS),
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

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/devpies/core/admin/pkg/handler"
	"github.com/devpies/core/admin/pkg/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devpies/core/pkg/log"
)

// Config contains application configuration with good defaults.
type config struct {
	Web struct {
		Address         string        `conf:"default:localhost:4000"`
		Debug           string        `conf:"default:localhost:6060"`
		Production      bool          `conf:"default:false"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		APIAddress      string        `conf:"default:http://localhost:4001"`
	}
	Stripe struct {
		Key    string `conf:"default:none"`
		Secret string `conf:"default:none"`
	}
}

var logPath = ".log/out.log"

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	var (
		cfg config
		err error
	)

	logger, Sync := log.NewLoggerOrPanic(logPath)
	defer Sync()

	if err = conf.Parse(os.Args[1:], "API", &cfg); err != nil {
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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	authService := service.NewAuthService(logger)
	authHandler := handler.NewAuth(logger, authService)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
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

package main

import (
	"context"
	"fmt"

	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/devpies/devpie-client-core/users/api/handlers"
	"github.com/devpies/devpie-client-core/users/api/listeners"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Port             string        `conf:"default::4000"`
			Debug            string        `conf:"default:localhost:6060"`
			Production       bool          `conf:"default:false"`
			ReadTimeout      time.Duration `conf:"default:5s"`
			WriteTimeout     time.Duration `conf:"default:5s"`
			ShutdownTimeout  time.Duration `conf:"default:5s"`
			CorsOrigins      string        `conf:"default:https://localhost:3000"`
			AuthDomain       string        `conf:"default:none"`
			AuthAudience     string        `conf:"default:none"`
			AuthM2MClient    string        `conf:"default:none"`
			AuthM2MSecret    string        `conf:"default:none"`
			AuthMAPIAudience string        `conf:"default:none"`
			SendgridAPIKey   string        `conf:"default:none,noprint"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		Nats struct {
			URL       string `conf:"default:nats://"`
			ClientID  string `conf:"default:client-id"`
			ClusterID string `conf:"default:cluster-id"`
		}
	}

	if err := conf.Parse(os.Args[1:], "API", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var help string
			help, err = conf.Usage("API", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(help)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	infolog := log.New(os.Stdout, fmt.Sprintf("%s:", cfg.Nats.ClientID), log.Lmsgprefix|log.Lmicroseconds|log.Lshortfile)

	infolog.Printf("main : Started")
	defer infolog.Println("main : Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}

	infolog.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Enabled Profiler

	go func() {
		log.Printf("main: Debug service listening on %s", cfg.Web.Debug)
		err = http.ListenAndServe(cfg.Web.Debug, nil)
		if err != nil {
			log.Printf("main: Debug service listening on %s", cfg.Web.Debug)
		}
	}()

	// =========================================================================
	// Start Database

	repo, rClose, err := database.NewRepository(database.Config{
		User:       cfg.DB.User,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer rClose()

	// =========================================================================
	// Start Listeners
	rand.New(rand.NewSource(time.Now().UnixNano()))
	clusterID := fmt.Sprintf("%s-%d", cfg.Nats.ClientID, rand.Int())
	queueGroup := fmt.Sprintf("%s-queue", cfg.Nats.ClientID)

	nats, eClose := events.NewClient(cfg.Nats.ClusterID, clusterID, cfg.Nats.URL)
	defer func() {
		cerr := eClose()
		if cerr != nil {
			err = cerr
		}
	}()

	go func(repo *database.Repository, nats *events.Client, infolog *log.Logger, queueGroup string) {
		l := listeners.NewListener(infolog, repo)
		l.RegisterAll(nats, queueGroup)
	}(repo, nats, infolog, queueGroup)

	// =========================================================================
	// Start API Service

	// Make a channel to listen for shutdown signal from the OS.
	shutdown := make(chan os.Signal, 1)

	api := http.Server{
		Addr: cfg.Web.Port,
		Handler: handlers.API(shutdown, repo, infolog, cfg.Web.CorsOrigins, cfg.Web.AuthAudience,
			cfg.Web.AuthDomain, cfg.Web.AuthMAPIAudience, cfg.Web.AuthM2MClient, cfg.Web.AuthM2MSecret,
			cfg.Web.SendgridAPIKey, nats),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)
	//
	//// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "listening and serving")

	case sig := <-shutdown:
		log.Println("main : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}

// Package billing provides the billing api.
package billing

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/devpies/saas-core/internal/billing/config"
	"github.com/devpies/saas-core/internal/billing/db"
	"github.com/devpies/saas-core/internal/billing/handler"
	"github.com/devpies/saas-core/internal/billing/repository"
	"github.com/devpies/saas-core/internal/billing/res"
	"github.com/devpies/saas-core/internal/billing/service"
	"github.com/devpies/saas-core/internal/billing/stripe"
	"github.com/devpies/saas-core/pkg/log"

	"go.uber.org/zap"
)

// Run contains the app setup.
func Run() error {
	var (
		logger  *zap.Logger
		logPath = "log/out.log"
		err     error
	)

	cfg, err := config.NewConfig()
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

	pg, Close, err := db.NewPostgresDatabase(logger, cfg)
	if err != nil {
		return err
	}
	defer Close()

	// Execute latest migration in production.
	if cfg.Web.Production {
		if err = res.MigrateUp(pg.URL.String()); err != nil {
			logger.Error("error connecting to user database", zap.Error(err))
			return err
		}
	}

	stripeClient := stripe.NewStripeClient(logger, cfg.Stripe.Key, cfg.Stripe.Secret)

	// Initialize 3-layered architecture.
	subscriptionRepository := repository.NewSubscriptionRepository(logger, pg)
	transactionRepository := repository.NewTransactionRepository(logger, pg)
	customerRepository := repository.NewCustomerRepository(logger, pg)

	subscriptionService := service.NewSubscriptionService(logger, stripeClient, subscriptionRepository)
	transactionService := service.NewTransactionService(logger, transactionRepository)
	customerService := service.NewCustomerService(logger, customerRepository)

	subscriptionHandler := handler.NewSubscriptionHandler(logger, subscriptionService, transactionService, customerService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, subscriptionHandler, cfg),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting project service on %s:%s", cfg.Web.Address, cfg.Web.Port))
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

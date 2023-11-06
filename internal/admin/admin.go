// Package admin provides a saas admin web app for administrators.
package admin

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devpies/saas-core/internal/admin/clients"
	"github.com/devpies/saas-core/internal/admin/config"
	"github.com/devpies/saas-core/internal/admin/db"
	"github.com/devpies/saas-core/internal/admin/handler"
	"github.com/devpies/saas-core/internal/admin/render"
	"github.com/devpies/saas-core/internal/admin/res"
	"github.com/devpies/saas-core/internal/admin/service"
	"github.com/devpies/saas-core/pkg/log"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

var session *scs.SessionManager

// Run contains the app setup.
func Run(staticFS embed.FS) error {
	var (
		cfg     config.Config
		logger  *zap.Logger
		logPath = "log/out.log"
		err     error
	)

	// Initialize static files.
	templateFS, err := fs.Sub(staticFS, "static/templates")
	if err != nil {
		logger.Error("error retrieving static templates", zap.Error(err))
		return err
	}
	assets, err := fs.Sub(staticFS, "static/assets")
	if err != nil {
		logger.Error("error retrieving static assets", zap.Error(err))
		return err
	}

	cfg, err = config.NewConfig()
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

	// Initialize admin database.
	database, Close, err := db.NewPostgresDatabase(logger, cfg)
	if err != nil {
		logger.Error("error connecting to admin database", zap.Error(err))
		return err
	}
	defer Close()

	// Execute latest migration.
	if cfg.Web.Production {
		if err = res.MigrateUp(database.URL.String()); err != nil {
			logger.Error("error connecting to admin database", zap.Error(err))
			return err
		}
	}

	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = postgresstore.New(database.DB.DB)

	ctx := context.Background()
	cognitoClient := clients.NewCognitoClient(ctx, cfg.Cognito.Region)
	billingClient := clients.NewHTTPBillingClient(
		logger,
		cfg.Billing.ServiceAddress,
		cfg.Billing.ServicePort,
		cognitoClient,
		cfg.Cognito.SharedUserPoolClientID,
		cfg.Cognito.SharedUserPoolID,
		cfg.Cognito.M2MClientKey,
		cfg.Cognito.M2MClientSecret,
	)
	registrationClient := clients.NewHTTPRegistrationClient(logger, cfg.Registration.ServiceAddress, cfg.Registration.ServicePort)
	tenantClient := clients.NewHTTPTenantClient(logger, cfg.Tenant.ServiceAddress, cfg.Tenant.ServicePort)

	// Initialize 3-layered architecture.
	authService := service.NewAuthService(logger, cfg.Cognito.Region, cfg.Cognito.UserPoolClientID, cfg.Cognito.UserPoolID, cognitoClient, session)
	registrationService := service.NewRegistrationService(logger, cfg.Cognito.SharedUserPoolID, cognitoClient, registrationClient)
	tenantService := service.NewTenantService(logger, tenantClient, billingClient)
	renderEngine := render.New(logger, cfg, templateFS, session)
	authHandler := handler.NewAuthHandler(logger, renderEngine, session, authService)
	webPageHandler := handler.NewWebPageHandler(logger, renderEngine, web.SetContextStatusCode)
	registrationHandler := handler.NewRegistrationHandler(logger, registrationService)
	tenantHandler := handler.NewTenantHandler(logger, tenantService)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, assets, cfg.Cognito.Region, cfg.Cognito.UserPoolID, authHandler, webPageHandler, registrationHandler, tenantHandler),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting admin app on %s:%s", cfg.Web.Address, cfg.Web.Port))
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

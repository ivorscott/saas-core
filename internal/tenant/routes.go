package tenant

import (
	"net/http"
	"os"

	"github.com/devpies/saas-core/internal/tenant/config"
	"github.com/devpies/saas-core/internal/tenant/handler"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Routes composes routes, middleware and handlers.
func Routes(
	log *zap.Logger,
	shutdown chan os.Signal,
	tenantHandler *handler.TenantHandler,
	authInfoHandler *handler.AuthInfoHandler,
	config config.Config,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://admin.devpie.local", "https://admin.devpie.io"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	middleware := []web.Middleware{
		mid.Logger(log),
		mid.Errors(log),
		mid.Auth(log, config.Cognito.Region, config.Cognito.UserPoolID),
		mid.Panics(log),
	}

	app := web.NewApp(mux, shutdown, log, middleware...)

	app.Handle(http.MethodGet, "/tenants", tenantHandler.FindAll)
	app.Handle(http.MethodGet, "/tenants/{id}", tenantHandler.FindOne)
	app.Handle(http.MethodPatch, "/tenants/{id}", tenantHandler.Update)
	app.Handle(http.MethodDelete, "/tenants/{id}", tenantHandler.Delete)
	app.Handle(http.MethodGet, "/tenants/auth-info", authInfoHandler.GetAuthInfo)

	return mux
}

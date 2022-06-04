package admin

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/devpies/saas-core/internal/admin/handler"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Routes composes routes, middleware and handlers.
func Routes(
	log *zap.Logger,
	shutdown chan os.Signal,
	assets fs.FS,
	authHandler *handler.AuthHandler,
	webPageHandler *handler.WebPageHandler,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(loadSession)

	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(assets))))

	app := web.NewApp(mux, shutdown, log, []web.Middleware{mid.Logger(log), mid.Errors(log), mid.Panics(log)}...)

	// Unauthenticated webpages.
	app.Handle(http.MethodGet, "/", withNoSession()(authHandler.Login))
	app.Handle(http.MethodGet, "/force-new-password", withPasswordChallengeSession()(authHandler.ForceNewPassword))
	app.Handle(http.MethodPost, "/secure-new-password", withNoSession()(authHandler.SetupNewUserWithSecurePassword))
	app.Handle(http.MethodPost, "/authenticate", withNoSession()(authHandler.AuthenticateCredentials))

	// Authenticated webpages.
	app.Handle(http.MethodGet, "/admin", withSession()(webPageHandler.Dashboard))
	app.Handle(http.MethodGet, "/admin/tenants", withSession()(webPageHandler.Tenants))
	app.Handle(http.MethodGet, "/admin/create-tenant", withSession()(webPageHandler.CreateTenant))
	app.Handle(http.MethodGet, "/admin/logout", withSession()(authHandler.Logout))
	app.Handle(http.MethodGet, "/*", withSession()(webPageHandler.E404))

	return mux
}

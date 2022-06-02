package adminclient

import (
	"github.com/devpies/core/internal/adminclient/handler"
	"github.com/devpies/core/pkg/web"
	"github.com/devpies/core/pkg/web/mid"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
	"os"

	"github.com/devpies/core/internal/adminclient/webpage"

	"github.com/go-chi/chi/v5"
)

// API composes routes, middleware and handlers.
func API(
	log *zap.Logger,
	shutdown chan os.Signal,
	assets fs.FS,
	page *webpage.WebPage,
	authHandler *handler.AuthHandler,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(loadSession)

	// Static webpages templates
	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(assets))))

	// Unauthenticated webpages.
	mux.Get("/", page.Login)
	mux.Get("/setup/new-password", page.ForceNewPassword)

	// Authenticated webpages.
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(withSession)
		mux.Get("/dashboard", page.Dashboard)
	})

	// Admin API endpoints.
	app := web.New(mux, shutdown, log, mid.Logger(log), mid.Errors(log), mid.Panics(log))
	app.Handle(http.MethodGet, "/api/logout", authHandler.Logout)
	app.Handle(http.MethodPost, "/api/authenticate", authHandler.AuthenticateCredentials)
	app.Handle(http.MethodPost, "/api/setup/new-user", authHandler.SetupNewUserWithSecurePassword)

	return mux
}

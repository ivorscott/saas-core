package admin

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/devpies/core/internal/admin/config"
	"github.com/devpies/core/internal/admin/webapp"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// API composes routes, middleware and handlers.
func API(
	log *zap.Logger,
	shutdown chan os.Signal,
	cfg config.Config,
	assets fs.FS,
	app *webapp.WebApp,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(loadSession)

	// Static webpages templates
	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(assets))))

	// Unauthenticated webpages.
	mux.Get("/", app.Login)
	mux.Get("/setup/new-password", app.ForceNewPassword)
	mux.Post("/authenticate", app.AuthenticateCredentials)
	mux.Post("/force-new-password", app.SetupNewUserWithSecurePassword)

	// Authenticated webpages.
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(withSession)
		mux.Get("/", app.Dashboard)
		mux.Get("/logout", app.Logout)
	})

	// Admin API endpoints.
	//app := web.New(mux, shutdown, log, mid.Logger(log), mid.Auth(log, cfg.Cognito.Region, cfg.Cognito.UserPoolClientID), mid.Errors(log), mid.Panics(log))

	return mux
}

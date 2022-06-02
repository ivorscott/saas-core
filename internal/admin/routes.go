package admin

import (
	"io/fs"
	"net/http"

	"github.com/devpies/core/internal/admin/handler"

	"github.com/go-chi/chi/v5"
)

// Routes composes routes, middleware and handlers.
func Routes(
	assets fs.FS,
	authHandler *handler.AuthHandler,
	webPageHandler *handler.WebPageHandler,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(loadSession)

	// Static webpages templates
	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(assets))))

	// Unauthenticated webpages.
	mux.Get("/", authHandler.Login)
	mux.Get("/setup/new-password", authHandler.ForceNewPassword)
	mux.Post("/authenticate", authHandler.AuthenticateCredentials)
	mux.Post("/force-new-password", authHandler.SetupNewUserWithSecurePassword)

	// Authenticated webpages.
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(withSession)
		mux.Get("/", webPageHandler.Dashboard)
		mux.Get("/logout", authHandler.Logout)
	})

	mux.Get("/*", webPageHandler.E404)

	return mux
}

package adminclient

import (
	"io/fs"
	"net/http"

	"github.com/devpies/core/internal/adminclient/webpage"

	"github.com/go-chi/chi/v5"
)

// API composes routes, middleware and handlers.
func API(
	assets fs.FS,
	page *webpage.WebPage,
) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", page.Login)
	mux.Get("/admin/dashboard", page.Dashboard)
	mux.Get("/setup/new-password", page.ForceNewPassword)
	mux.Get("/logout", page.Logout)

	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(assets))))

	return mux
}

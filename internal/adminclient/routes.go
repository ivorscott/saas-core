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
	mux.Use(SessionLoad)

	mux.Get("/", page.Login)
	mux.Get("/setup/new-password", page.ForceNewPassword)
	mux.Get("/logout", page.Logout)
	// FIXME: It would be cleaner to do the whole server-side auth flow in adminclient, get the token and make a session
	// The session would be easier to create this way without the back and forth communication between api and client
	// Then adminapi can require the token for every endpoint request.
	mux.Get("/sess/{id}", page.CreateSession)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(withAuth)
		mux.Get("/dashboard", page.Dashboard)
	})

	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(assets))))

	return mux
}

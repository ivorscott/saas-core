package main

import (
	"embed"
	"github.com/devpies/core/admin/pkg/webapp"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
	"os"
)

// API composes routes, middleware and handlers.
func API(
	shutdown chan os.Signal,
	logger *zap.Logger,
	staticFS embed.FS,
	app *webapp.WebApp,
) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", app.Login)
	mux.Get("/logout", app.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(withAuth)
		mux.Get("/", app.Dashboard)
	})

	assets, err := fs.Sub(staticFS, "static/assets")
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(assets))))

	return mux
}

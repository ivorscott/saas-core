package main

import (
	"embed"
	"github.com/devpies/core/admin/pkg/render"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
	"os"
)

// API configures the application routes, middleware and handlers.
func API(
	shutdown chan os.Signal,
	logger *zap.Logger,
	render *render.Render,
	contentFS embed.FS,
) http.Handler {
	mux := chi.NewRouter()

	swagger, _ := fs.Sub(content, "te/swagger-ui")

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", Login)
	mux.Get("/logout", Logout)
	mux.Get("/ws", WsEndpoint)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(withAuth)
		mux.Get("/", Dashboard)
	})

	return mux
}

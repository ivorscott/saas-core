package main

import (
	"embed"
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
	content embed.FS,
) http.Handler {
	mux := chi.NewRouter()

	js, _ := fs.Sub(content, "static/js")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(js))))

	mux.Get("/", Login)
	mux.Get("/logout", Logout)
	mux.Get("/ws", WsEndpoint)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(withAuth)
		mux.Get("/", Dashboard)
	})

	return mux
}

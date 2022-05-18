package main

import (
	"go.uber.org/zap"
	"net/http"
	"os"

	"github.com/devpies/core/admin/pkg/handler"
	"github.com/go-chi/chi/v5"
)

// API configures the application routes, middleware and handlers.
func API(
	shutdown chan os.Signal,
	logger *zap.Logger,
	authHandler *handler.AuthHandler,
) http.Handler {
	mux := chi.NewRouter()

	mux.Route("/api", func(mux chi.Router) {
		mux.Use(withAuth)
		mux.Get("/auth", authHandler.AuthenticateCredentials)
	})

	return mux
}

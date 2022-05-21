package main

import (
	"github.com/devpies/core/pkg/web"
	"github.com/devpies/core/pkg/web/mid"
	"github.com/go-chi/cors"

	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/devpies/core/internal/adminapi/handler"
)

// API configures the application routes, middleware and handlers.
func API(
	shutdown chan os.Signal,
	log *zap.Logger,
	authHandler *handler.AuthHandler,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	app := web.New(mux, shutdown, log, mid.Logger(log), mid.Errors(log), mid.Panics(log))
	app.Handle(http.MethodPost, "/api/authenticate", authHandler.AuthenticateCredentials)
	app.Handle(http.MethodPost, "/api/setup/new-user", authHandler.SetupNewUserWithSecurePassword)

	return mux
}

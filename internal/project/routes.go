package project

import (
	"net/http"
	"os"

	"github.com/devpies/saas-core/internal/project/config"
	"github.com/devpies/saas-core/internal/project/handler"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Routes composes routes, middleware and handlers.
func Routes(
	log *zap.Logger,
	shutdown chan os.Signal,
	taskHandler *handler.TaskHandler,
	columnHandler *handler.ColumnHandler,
	projectHandler *handler.ProjectHandler,
	config config.Config,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "https://devpie.io"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	middleware := []web.Middleware{
		mid.Logger(log),
		mid.Errors(log),
		//mid.Auth(log, config.Cognito.Region, config.Cognito.UserPoolClientID),
		func(h web.Handler) web.Handler {
			// fake auth middleware
			// set fake tenant id
			return h
		},
		mid.Panics(log),
	}

	app := web.NewApp(mux, shutdown, log, middleware...)

	app.Handle(http.MethodGet, "/", projectHandler.List)
	app.Handle(http.MethodPost, "/", projectHandler.Create)
	app.Handle(http.MethodGet, "/{pid}", projectHandler.Retrieve)
	app.Handle(http.MethodPatch, "/{pid}", projectHandler.Update)
	app.Handle(http.MethodDelete, "/{pid}", projectHandler.Delete)
	app.Handle(http.MethodGet, "/{pid}/columns", columnHandler.List)
	app.Handle(http.MethodGet, "/{pid}/tasks", taskHandler.List)
	app.Handle(http.MethodPost, "/{pid}/columns/{cid}/tasks", taskHandler.Create)
	app.Handle(http.MethodPatch, "/tasks/{tid}", taskHandler.Update)
	app.Handle(http.MethodPatch, "/tasks/{tid}/move", taskHandler.Move)
	app.Handle(http.MethodDelete, "/columns/{cid}/tasks/{tid}", taskHandler.Delete)

	return mux
}

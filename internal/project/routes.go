package project

import (
	"github.com/go-chi/cors"
	"net/http"
	"os"

	"github.com/devpies/saas-core/internal/project/config"
	"github.com/devpies/saas-core/internal/project/handler"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"

	"github.com/go-chi/chi/v5"
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
		AllowedOrigins:   []string{"https://devpie.local:3000", "https://devpie.io"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "BasePath"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	middleware := []web.Middleware{
		mid.Logger(log),
		mid.Errors(log),
		mid.Auth(log, config.Cognito.Region, config.Cognito.UserPoolID),
		mid.Panics(log),
	}

	app := web.NewApp(mux, shutdown, log, middleware...)

	app.Handle(http.MethodGet, "/projects", projectHandler.List)
	app.Handle(http.MethodPost, "/projects", projectHandler.Create)
	app.Handle(http.MethodGet, "/projects/{pid}", projectHandler.Retrieve)
	app.Handle(http.MethodPatch, "/projects/{pid}", projectHandler.Update)
	app.Handle(http.MethodDelete, "/projects/{pid}", projectHandler.Delete)
	app.Handle(http.MethodGet, "/projects/{pid}/columns", columnHandler.List)
	app.Handle(http.MethodGet, "/projects/{pid}/tasks", taskHandler.List)
	app.Handle(http.MethodPost, "/projects/{pid}/columns/{cid}/tasks", taskHandler.Create)
	app.Handle(http.MethodPatch, "/projects/tasks/{tid}", taskHandler.Update)
	app.Handle(http.MethodPatch, "/projects/tasks/{tid}/move", taskHandler.Move)
	app.Handle(http.MethodDelete, "/projects/columns/{cid}/tasks/{tid}", taskHandler.Delete)

	return app
}

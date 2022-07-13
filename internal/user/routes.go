package user

import (
	"github.com/devpies/saas-core/internal/user/handler"
	"net/http"
	"os"

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
	region string,
	sharedUserPoolID string,
	userHandler *handler.UserHandler,
	inviteHandler *handler.InviteHandler,
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
		mid.Auth(log, region, sharedUserPoolID),
		mid.Panics(log),
	}

	app := web.NewApp(mux, shutdown, log, middleware...)

	app.Handle(http.MethodPost, "/users", userHandler.Create)
	app.Handle(http.MethodGet, "/users", userHandler.List)
	app.Handle(http.MethodGet, "/users/me", userHandler.RetrieveMe)
	app.Handle(http.MethodDelete, "/users/{uid}", userHandler.RemoveUser)
	app.Handle(http.MethodGet, "/users/available-seats", userHandler.SeatsAvailable)
	app.Handle(http.MethodGet, "/users/invites", inviteHandler.RetrieveInvites)
	app.Handle(http.MethodPost, "/users/invites", inviteHandler.CreateInvite)
	app.Handle(http.MethodPatch, "/users/invites/{iid}", inviteHandler.UpdateInvite)

	return app
}

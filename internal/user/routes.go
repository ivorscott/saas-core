package user

import (
	"github.com/devpies/saas-core/internal/user/handler"
	"net/http"
	"os"

	"github.com/devpies/saas-core/internal/user/config"
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
	userHandler *handler.UserHandler,
	teamHandler *handler.TeamHandler,
	membershipHandler *handler.MembershipHandler,
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
		//mid.Auth(log, config.Cognito.Region, config.Cognito.SharedUserPoolClientID),
		func(h web.Handler) web.Handler {
			// fake auth middleware
			// set fake tenant id
			return h
		},
		mid.Panics(log),
	}

	app := web.NewApp(mux, shutdown, log, middleware...)

	app.Handle(http.MethodPost, "/", userHandler.Create)
	app.Handle(http.MethodGet, "/me", userHandler.RetrieveMe)
	app.Handle(http.MethodPost, "/teams", teamHandler.Create)
	app.Handle(http.MethodPost, "/teams/{tid}/project/{pid}", teamHandler.AssignExistingTeam)
	app.Handle(http.MethodPost, "/teams/{tid}/leave", teamHandler.LeaveTeam)
	app.Handle(http.MethodGet, "/teams", teamHandler.List)
	app.Handle(http.MethodGet, "/teams/{tid}", teamHandler.Retrieve)
	app.Handle(http.MethodPost, "/teams/{tid}/invites", teamHandler.CreateInvite)
	app.Handle(http.MethodGet, "/teams/invites", teamHandler.RetrieveInvites)
	app.Handle(http.MethodGet, "/teams/{tid}/members", membershipHandler.RetrieveMemberships)
	app.Handle(http.MethodPatch, "/teams/{tid}/invites/{iid}", teamHandler.UpdateInvite)

	return mux
}

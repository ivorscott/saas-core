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
		mid.Auth(log, config.Cognito.Region, config.Cognito.SharedUserPoolClientID),
		mid.Panics(log),
	}

	app := web.NewApp(mux, shutdown, log, middleware...)

	//app.Handle(http.MethodPost, "/users/", userHandler.Create)
	app.Handle(http.MethodGet, "/users/me", userHandler.RetrieveMe)
	app.Handle(http.MethodPost, "/users/teams", teamHandler.CreateTeamForProject)
	app.Handle(http.MethodPost, "/users/teams/{tid}/project/{pid}", teamHandler.AssignExistingTeam)
	app.Handle(http.MethodPost, "/users/teams/{tid}/leave", teamHandler.LeaveTeam)
	app.Handle(http.MethodGet, "/users/teams", teamHandler.List)
	app.Handle(http.MethodGet, "/users/teams/{tid}", teamHandler.Retrieve)
	app.Handle(http.MethodPost, "/users/teams/{tid}/invites", teamHandler.CreateInvite)
	app.Handle(http.MethodGet, "/users/teams/invites", teamHandler.RetrieveInvites)
	app.Handle(http.MethodGet, "/users/teams/{tid}/members", membershipHandler.RetrieveMemberships)
	app.Handle(http.MethodPatch, "/users/teams/{tid}/invites/{iid}", teamHandler.UpdateInvite)

	return mux
}

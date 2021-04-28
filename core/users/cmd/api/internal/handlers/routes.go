package handlers

import (
	"github.com/ivorscott/devpie-client-core/users/internal/mid"
	"github.com/ivorscott/devpie-client-core/users/internal/platform/database"
	"github.com/ivorscott/devpie-client-core/users/internal/platform/web"
	"log"
	"net/http"
	"os"
)

func API(shutdown chan os.Signal, repo *database.Repository, log *log.Logger, origins string,
	Auth0Audience, Auth0Domain, Auth0MAPIAudience, Auth0M2MClient, Auth0M2MSecret, SendgridAPIKey string) http.Handler {

	auth0 := &mid.Auth0{
		Audience:     Auth0Audience,
		Domain:       Auth0Domain,
		MAPIAudience: Auth0MAPIAudience,
		M2MClient:    Auth0M2MClient,
		M2MSecret:    Auth0M2MSecret,
	}

	app := web.NewApp(shutdown, log, mid.Logger(log), auth0.Authenticate(), mid.Errors(log), mid.Panics(log))

	h := HealthCheck{repo: repo}

	app.Handle(http.MethodGet, "/api/v1/health", h.Health)

	u := Users{repo: repo, log: log, auth0: auth0, origins: origins, sendgridAPIKey: SendgridAPIKey}
	tm := Team{repo: repo, log: log, auth0: auth0, origins: origins, sendgridAPIKey: SendgridAPIKey}

	app.Handle(http.MethodPost, "/api/v1/users", u.Create)
	app.Handle(http.MethodGet, "/api/v1/users/me", u.RetrieveMe)

	app.Handle(http.MethodPost, "/api/v1/users/teams", tm.Create)
	app.Handle(http.MethodGet, "/api/v1/users/teams/{tid}", tm.Retrieve)
	app.Handle(http.MethodGet, "/api/v1/users/teams/invites", tm.RetrieveInvites)
	app.Handle(http.MethodPost, "/api/v1/users/teams/{tid}/invite", tm.Invite)

	return Cors(origins).Handler(app)
}

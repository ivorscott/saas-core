package handlers

import (
	"log"
	"net/http"
	"os"

	mid "github.com/devpies/devpie-client-core/users/api/middleware"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/devpies/devpie-client-events/go/events"
)

func API(shutdown chan os.Signal, repo *database.Repository, log *log.Logger, origins string,
	auth0Audience, auth0Domain, auth0MAPIAudience, auth0M2MClient, auth0M2MSecret, sendgridKey string, nats *events.Client) http.Handler {

	a0 := &auth0.Auth0{
		Repo:         repo,
		Domain:       auth0Domain,
		Audience:     auth0Audience,
		M2MSecret:    auth0M2MSecret,
		M2MClient:    auth0M2MClient,
		MAPIAudience: auth0MAPIAudience,
	}

	app := web.NewApp(shutdown, log, mid.Logger(log), a0.Authenticate(), mid.Errors(log), mid.Panics(log))

	h := HealthCheck{repo: repo}

	app.Handle(http.MethodGet, "/api/v1/health", h.Health)

	u := Users{repo, log, a0, origins}
	tm := Team{repo, log, a0, nats, origins, sendgridKey}

	app.Handle(http.MethodPost, "/api/v1/users", u.Create)
	app.Handle(http.MethodGet, "/api/v1/users/me", u.RetrieveMe)

	app.Handle(http.MethodPost, "/api/v1/users/teams", tm.Create)
	//app.Handle(http.MethodGet, "/api/v1/users/teams", tm.List)
	app.Handle(http.MethodGet, "/api/v1/users/teams/{tid}", tm.Retrieve)
	app.Handle(http.MethodPost, "/api/v1/users/teams/{tid}/invites", tm.CreateInvite)
	app.Handle(http.MethodGet, "/api/v1/users/teams/invites", tm.RetrieveInvites)
	app.Handle(http.MethodPatch, "/api/v1/users/teams/{tid}/invites/{iid}", tm.UpdateInvite)

	return Cors(origins).Handler(app)
}

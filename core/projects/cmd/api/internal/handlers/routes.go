package handlers

import (
	"fmt"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/devpies/devpie-client-core/projects/cmd/api/internal/listeners"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/devpies/devpie-client-core/projects/internal/mid"
	"github.com/devpies/devpie-client-core/projects/internal/platform/database"
	"github.com/devpies/devpie-client-core/projects/internal/platform/web"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func API(shutdown chan os.Signal, repo *database.Repository, log *log.Logger, origins string,
	Auth0Audience, Auth0Domain, Auth0MAPIAudience, Auth0M2MClient, Auth0M2MSecret, NatsURL,
	NatsClientId, NatsClusterId string) http.Handler {

	clusterId := fmt.Sprintf("%s-%d", NatsClientId, rand.Int())
	queueGroup := fmt.Sprintf("%s-queue", NatsClientId)

	auth0 := &mid.Auth0{
		Audience:     Auth0Audience,
		Domain:       Auth0Domain,
		MAPIAudience: Auth0MAPIAudience,
		M2MClient:    Auth0M2MClient,
		M2MSecret:    Auth0M2MSecret,
	}

	app := web.NewApp(shutdown, log, mid.Logger(log), auth0.Authenticate(), mid.Errors(log), mid.Panics(log))

	nats, close := events.NewClient(NatsClusterId, clusterId, NatsURL)
	defer close()

	l := listeners.NewListeners(log, repo)

	h := HealthCheck{repo: repo}

	app.Handle(http.MethodGet, "/api/v1/health", h.Health)

	t := Tasks{repo: repo, log: log, auth0: auth0}
	c := Columns{repo: repo, log: log, auth0: auth0}
	p := Projects{repo: repo, log: log, auth0: auth0}

	app.Handle(http.MethodGet, "/api/v1/projects", p.List)
	app.Handle(http.MethodPost, "/api/v1/projects", p.Create)
	app.Handle(http.MethodGet, "/api/v1/projects/{pid}", p.Retrieve)
	app.Handle(http.MethodPut, "/api/v1/projects/{pid}", p.Update)
	app.Handle(http.MethodDelete, "/api/v1/projects/{pid}", p.Delete)
	app.Handle(http.MethodGet, "/api/v1/projects/{pid}/columns", c.List)
	app.Handle(http.MethodGet, "/api/v1/projects/{pid}/tasks", t.List)
	app.Handle(http.MethodPost, "/api/v1/projects/{pid}/columns/{cid}/tasks", t.Create)
	app.Handle(http.MethodPatch, "/api/v1/projects/tasks/{tid}", t.Update)
	app.Handle(http.MethodPatch, "/api/v1/projects/tasks/{tid}/move", t.Move)
	app.Handle(http.MethodDelete, "/api/v1/projects/columns/{cid}/tasks/{tid}", t.Delete)

	l.RegisterAll(nats, queueGroup)

	return Cors(origins).Handler(app)
}

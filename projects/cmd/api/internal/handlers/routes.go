package handlers

import (
	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)


func API(shutdown chan os.Signal, repo *database.Repository, log *log.Logger, origins []string,
	Auth0Audience, Auth0Domain string) http.Handler {

	auth0 := &mid.Auth0{
		Audience:     Auth0Audience,
		Domain:       Auth0Domain,
	}

	app := web.NewApp(shutdown, log, mid.Logger(log), auth0.Authenticate(), mid.Errors(log), mid.Panics(log))

	cr := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedHeaders:   []string{"Authorization", "Cache-Control", "Content-Type", "Strict-Transport-Security"},
		AllowedMethods:   []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
		AllowCredentials: true,
	})

	h := HealthCheck{repo: repo}

	app.Handle(http.MethodGet, "/v1/health", h.Health)

	t := Tasks{repo: repo, log: log, auth0: auth0}
	c := Columns{repo: repo, log: log, auth0: auth0}
	p := Projects{repo: repo, log: log, auth0: auth0}

	app.Handle(http.MethodGet, "/v1/projects", p.List)

	app.Handle(http.MethodPost, "/v1/projects", p.Create)
	app.Handle(http.MethodGet, "/v1/projects/{pid}", p.Retrieve)
	app.Handle(http.MethodPut, "/v1/projects/{pid}", p.Update)
	app.Handle(http.MethodDelete, "/v1/projects/{pid}", p.Delete)
	app.Handle(http.MethodGet, "/v1/projects/{pid}/columns", c.List)
	app.Handle(http.MethodGet, "/v1/projects/{pid}/tasks", t.List)
	app.Handle(http.MethodPost, "/v1/projects/{pid}/columns/{cid}/tasks", t.Create)
	app.Handle(http.MethodPatch, "/v1/projects/{pid}/tasks/{tid}", t.Update)
	app.Handle(http.MethodPatch, "/v1/projects/{pid}/tasks/{tid}/move", t.Move)
	app.Handle(http.MethodDelete, "/v1/projects/{pid}/columns/{cid}/tasks/{tid}", t.Delete)

	return cr.Handler(app)
}

package handlers

import (
	"net/http"

	"github.com/devpies/devpie-client-core/projects/internal/platform/database"
	"github.com/devpies/devpie-client-core/projects/internal/platform/web"
)

type HealthCheck struct {
	repo *database.Repository
}

func (c *HealthCheck) Health(w http.ResponseWriter, r *http.Request) error {
	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(r.Context(), c.repo.DB); err != nil {
		health.Status = "db not ready"
		return web.Respond(r.Context(), w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return web.Respond(r.Context(), w, health, http.StatusOK)
}

package handlers

import (
	"net/http"

	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/web"
)

// HealthCheck defines the service's health check mechanism
type HealthCheck struct {
	repo database.Storer
}

// Health returns the service's health status
func (c *HealthCheck) Health(w http.ResponseWriter, r *http.Request) error {
	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(r.Context(), c.repo); err != nil {
		health.Status = "db not ready"
		return web.Respond(r.Context(), w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return web.Respond(r.Context(), w, health, http.StatusOK)
}

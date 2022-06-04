package handler

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type setStatusCodeFunc func(ctx context.Context, statusCode int)

// WebPageHandler renders various webpages required for SaaS administration.
type WebPageHandler struct {
	logger        *zap.Logger
	render        renderer
	setStatusCode setStatusCodeFunc
}

// NewWebPageHandler returns a new webpage handler.
func NewWebPageHandler(
	logger *zap.Logger,
	renderEngine renderer,
	setStatus setStatusCodeFunc,
) *WebPageHandler {
	return &WebPageHandler{
		logger:        logger,
		render:        renderEngine,
		setStatusCode: setStatus,
	}
}

// Dashboard displays a useful dashboard for users.
func (page *WebPageHandler) Dashboard(w http.ResponseWriter, r *http.Request) error {
	return page.render.Template(w, r, "dashboard", nil)
}

// Tenants displays a table list for tenants.
func (page *WebPageHandler) Tenants(w http.ResponseWriter, r *http.Request) error {
	return page.render.Template(w, r, "tenants", nil)
}

// CreateTenant displays a tenant registration form.
func (page *WebPageHandler) CreateTenant(w http.ResponseWriter, r *http.Request) error {
	return page.render.Template(w, r, "create-tenant", nil)
}

// E404 displays a 404 error page.
func (page *WebPageHandler) E404(w http.ResponseWriter, r *http.Request) error {
	err := page.render.Template(w, r, "404", nil)
	page.setStatusCode(r.Context(), http.StatusNotFound)
	return err
}

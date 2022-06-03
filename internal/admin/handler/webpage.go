package handler

import (
	"net/http"

	"github.com/devpies/core/internal/admin/config"
	"github.com/devpies/core/internal/admin/render"
	"github.com/devpies/core/pkg/web"

	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

// WebPageHandler renders various webpages required for SaaS administration.
type WebPageHandler struct {
	logger  *zap.Logger
	config  config.Config
	render  *render.Render
	service authService
	session *scs.SessionManager
}

// NewWebPageHandler returns a new webpage handler.
func NewWebPageHandler(logger *zap.Logger, config config.Config, renderEngine *render.Render, service authService, session *scs.SessionManager) *WebPageHandler {
	return &WebPageHandler{
		logger:  logger,
		config:  config,
		render:  renderEngine,
		service: service,
		session: session,
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
	web.SetContextStatusCode(r.Context(), http.StatusNotFound)
	return err
}

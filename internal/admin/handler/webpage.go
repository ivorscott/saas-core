package handler

import (
	"net/http"

	"github.com/devpies/core/internal/admin/config"
	"github.com/devpies/core/internal/admin/render"

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

// E404 displays a 404 error page.
func (page *WebPageHandler) E404(w http.ResponseWriter, r *http.Request) error {
	return page.render.Template(w, r, "404", nil)
}

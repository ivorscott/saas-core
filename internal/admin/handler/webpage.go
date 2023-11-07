package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/render"

	"github.com/google/uuid"
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

// DashboardPage displays a useful dashboard for users.
func (page *WebPageHandler) DashboardPage(w http.ResponseWriter, r *http.Request) error {
	return page.render.Template(w, r, "dashboard", nil)
}

// TenantsPage displays a table list for tenants.
func (page *WebPageHandler) TenantsPage(w http.ResponseWriter, r *http.Request) error {
	return page.render.Template(w, r, "tenants", nil)
}

// TenantPage displays a specific tenant's details.
func (page *WebPageHandler) TenantPage(w http.ResponseWriter, r *http.Request) error {
	var data = make(map[string]string)
	//data["UserID"] = chi.URLParam(r, "id")
	return page.render.Template(w, r, "tenant-detail", &render.TemplateData{StringMap: data})
}

// CreateTenantPage displays a tenant registration form.
func (page *WebPageHandler) CreateTenantPage(w http.ResponseWriter, r *http.Request) error {
	var data = make(map[string]string)
	data["UserID"] = uuid.New().String()
	return page.render.Template(w, r, "create-tenant", &render.TemplateData{StringMap: data})
}

// E404Page displays a 404 error page.
func (page *WebPageHandler) E404Page(w http.ResponseWriter, r *http.Request) error {
	err := page.render.Template(w, r, "404", nil)
	page.setStatusCode(r.Context(), http.StatusNotFound)
	return err
}

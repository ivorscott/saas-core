// Package handler manages the presentation layer for handling incoming requests.
package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/tenant/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type tenantService interface {
	FindOne(ctx context.Context, tenantID string) (model.Tenant, error)
	FindAll(ctx context.Context) ([]model.Tenant, error)
	Update(ctx context.Context, id string, tenant model.UpdateTenant) error
	Delete(ctx context.Context, tenantID string) error
}

// TenantHandler handles tenant requests.
type TenantHandler struct {
	logger  *zap.Logger
	service tenantService
}

// NewTenantHandler returns a new TenantHandler.
func NewTenantHandler(logger *zap.Logger, service tenantService) *TenantHandler {
	return &TenantHandler{
		logger:  logger,
		service: service,
	}
}

// FindAll handles a search for all onboarded tenants.
func (th *TenantHandler) FindAll(w http.ResponseWriter, r *http.Request) error {
	tenants, err := th.service.FindAll(r.Context())
	if err != nil {
		return web.NewRequestError(err, http.StatusNotFound)
	}
	return web.Respond(r.Context(), w, tenants, http.StatusOK)
}

// FindOne handles a search for a single tenant.
func (th *TenantHandler) FindOne(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	tenant, err := th.service.FindOne(r.Context(), id)
	if err != nil {
		return web.NewRequestError(err, http.StatusNotFound)
	}
	return web.Respond(r.Context(), w, tenant, http.StatusOK)
}

// Update handles a tenant update request.
func (th *TenantHandler) Update(w http.ResponseWriter, r *http.Request) error {
	var (
		update model.UpdateTenant
		err    error
	)
	id := chi.URLParam(r, "id")
	err = web.Decode(r, &update)
	if err != nil {
		return err
	}
	err = th.service.Update(r.Context(), id, update)
	if err != nil {
		return web.NewRequestError(err, http.StatusNotFound)
	}
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// Delete handles a tenant deletion request.
func (th *TenantHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	err := th.service.Delete(r.Context(), id)
	if err != nil {
		return web.NewRequestError(err, http.StatusNotFound)
	}
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/tenant/model"

	"go.uber.org/zap"
)

type tenantService interface {
	AddConfiguration(ctx context.Context, tenantConfig model.NewTenantConfig) error
	GetAuthInfo(ctx context.Context, referer string) (model.AuthInfo, error)
	FindOne(ctx context.Context, tenantID string) (model.Tenant, error)
	FindAll(ctx context.Context) ([]model.Tenant, error)
	Update(ctx context.Context, tenant model.UpdateTenant) error
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

// GetAuthInfo handles a request for tenant authentication information.
func (th *TenantHandler) GetAuthInfo(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// FindAll handles a search for all onboarded tenants.
func (th *TenantHandler) FindAll(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// FindOne handles a search for a single tenant.
func (th *TenantHandler) FindOne(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Update handles a tenant update request.
func (th *TenantHandler) Update(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Delete handles a tenant deletion request.
func (th *TenantHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	return nil
}

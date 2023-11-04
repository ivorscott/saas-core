package handler

import (
	"context"
	"github.com/devpies/saas-core/internal/admin/model"
	"net/http"

	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type tenantService interface {
	ListTenants(ctx context.Context) ([]model.Tenant, int, error)
}

// TenantHandler handles the tenant requests.
type TenantHandler struct {
	logger  *zap.Logger
	service tenantService
}

// NewTenantHandler returns a new tenant handler.
func NewTenantHandler(
	logger *zap.Logger,
	service tenantService,
) *TenantHandler {
	return &TenantHandler{
		logger:  logger,
		service: service,
	}
}

// ListTenants lists all tenants.
func (th *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) error {
	tenants, status, err := th.service.ListTenants(r.Context())

	if err != nil {
		switch status {
		case http.StatusBadRequest:
			return web.NewRequestError(err, http.StatusBadRequest)
		case http.StatusNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case http.StatusUnauthorized:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			th.logger.Error("unexpected error", zap.Error(err), zap.Int("status", status))
			return err
		}
	}

	return web.Respond(r.Context(), w, tenants, http.StatusOK)
}

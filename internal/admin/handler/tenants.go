package handler

import (
	"net/http"

	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

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

// CancelSubscription cancels a stripe subscription for the customer.
func (th *TenantHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) error {
	var (
		subID = chi.URLParam(r, "subID")
		err   error
	)
	_, err = th.service.CancelSubscription(r.Context(), subID)
	if err != nil {
		return err
	}
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// RefundUser refunds the stripe customer.
func (th *TenantHandler) RefundUser(w http.ResponseWriter, r *http.Request) error {
	var (
		subID = chi.URLParam(r, "subID")
		err   error
	)
	_, err = th.service.RefundUser(r.Context(), subID)
	if err != nil {
		return err
	}
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

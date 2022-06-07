package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type registrationService interface {
	RegisterTenant(ctx context.Context, tenant model.NewTenant) (*web.ErrorResponse, int, error)
}

// RegistrationHandler handles the new tenant request from the admin app.
type RegistrationHandler struct {
	logger              *zap.Logger
	registrationService registrationService
}

// NewRegistrationHandler returns a new registration handler.
func NewRegistrationHandler(
	logger *zap.Logger,
	registrationService registrationService,
) *RegistrationHandler {
	return &RegistrationHandler{
		logger:              logger,
		registrationService: registrationService,
	}
}

// ProcessRegistration submits a new tenant to the registration service and responds with nil on success.
func (reg *RegistrationHandler) ProcessRegistration(w http.ResponseWriter, r *http.Request) error {
	var (
		payload model.NewTenant
		err     error
	)

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	webErrResp, status, err := reg.registrationService.RegisterTenant(r.Context(), payload)
	if err != nil {
		reg.logger.Info("registration client request failed", zap.Error(err))
		return err
	}
	if webErrResp != nil {
		switch status {
		case http.StatusBadRequest:
			return web.Respond(r.Context(), w, webErrResp, http.StatusBadRequest)
		default:
			return err
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

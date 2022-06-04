package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type registrationService interface {
	RegisterTenant(ctx context.Context, tenant model.NewTenant) error
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

// ProcessRegistration submits a new tenant to the registration service.
func (reg *RegistrationHandler) ProcessRegistration(w http.ResponseWriter, r *http.Request) error {
	var (
		payload model.NewTenant
		err     error
	)

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	// TODO: Check if tenant exists first.

	err = reg.registrationService.RegisterTenant(r.Context(), payload)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

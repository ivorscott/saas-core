package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/registration/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type registrationService interface {
	CreateRegistration(ctx context.Context, id string, tenant model.NewTenant) error
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

// RegisterTenant registers the new tenant.
func (reg *RegistrationHandler) RegisterTenant(w http.ResponseWriter, r *http.Request) error {
	var (
		payload model.NewTenant
		err     error
	)

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return web.NewShutdownError(err.Error())
	}
	err = reg.registrationService.CreateRegistration(r.Context(), id.String(), payload)
	if err != nil {
		reg.logger.Info("event publishing failed", zap.Error(err))
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

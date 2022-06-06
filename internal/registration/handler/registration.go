package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devpies/saas-core/internal/registration/model"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type registrationService interface {
	PublishTenantMessages(ctx context.Context, tenant model.NewTenant) error
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

	reg.logger.Info(fmt.Sprintf("%v", payload))

	// publish messages

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"go.uber.org/zap"
)

type registrationService interface {
	RegisterTenant(ctx context.Context, tenant model.NewTenant) (int, error)
	ResendTemporaryPassword(ctx context.Context, username string) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

// RegistrationHandler handles the new tenant request from the admin app.
type RegistrationHandler struct {
	logger  *zap.Logger
	service registrationService
}

// NewRegistrationHandler returns a new registration handler.
func NewRegistrationHandler(
	logger *zap.Logger,
	service registrationService,
) *RegistrationHandler {
	return &RegistrationHandler{
		logger:  logger,
		service: service,
	}
}

// ProcessRegistration submits a new tenant to the registration service and responds with nil on success.
func (rh *RegistrationHandler) ProcessRegistration(w http.ResponseWriter, r *http.Request) error {
	var (
		payload model.NewTenant
		err     error
	)

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	status, err := rh.service.RegisterTenant(r.Context(), payload)

	if err != nil {
		switch status {
		case http.StatusBadRequest:
			return web.NewRequestError(err, http.StatusBadRequest)
		case http.StatusNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case http.StatusUnauthorized:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			rh.logger.Error("unexpected error", zap.Error(err), zap.Int("status", status))
			return err
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// ResendOTP resends the one time password.
func (rh *RegistrationHandler) ResendOTP(w http.ResponseWriter, r *http.Request) error {
	var (
		payload struct {
			Username string `json:"username"`
		}
		err error
	)

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	if _, err = rh.service.ResendTemporaryPassword(r.Context(), payload.Username); err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

package handler

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type registrationService interface {
	RegisterTenant(ctx context.Context, tenant model.NewTenant) (*web.ErrorResponse, int, error)
	ResendTemporaryPassword(ctx context.Context, username string) (*cognitoidentityprovider.AdminCreateUserOutput, error)
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

func (reg *RegistrationHandler) ResendOTP(w http.ResponseWriter, r *http.Request) error {
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

	if _, err = reg.registrationService.ResendTemporaryPassword(r.Context(), payload.Username); err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

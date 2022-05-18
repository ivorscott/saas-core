package handler

import (
	"context"
	"github.com/devpies/core/admin/pkg/model"
	"github.com/devpies/core/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

type authService interface {
	Authenticate(ctx context.Context, email, password string) error
}

type AuthHandler struct {
	logger  *zap.Logger
	service authService
}

func NewAuth(logger *zap.Logger, service authService) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: service,
	}
}

// AuthenticateCredentials handles email and password values from the admin login form.
func (a *AuthHandler) AuthenticateCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		c   model.AuthCredentials
		err error
	)

	err = web.Decode(r, &c)
	if err != nil {
		a.logger.Error("", zap.Error(err))
	}
	
	err = a.service.Authenticate(r.Context(), c.Email, c.Password)
	if err != nil {
		a.logger.Error("authenticate failed", zap.Error(err))
	}

	err = web.Respond(r.Context(), w, nil, http.StatusOK)
	if err != nil {
		a.logger.Error("", zap.Error(err))
	}
}

package handler

import (
	"context"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/devpies/core/admin/pkg/model"
	"github.com/devpies/core/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

type authService interface {
	Authenticate(ctx context.Context, credentials model.AuthCredentials) (*cip.AdminInitiateAuthOutput, error)
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
		credentials model.AuthCredentials
		err         error
	)

	err = web.Decode(r, &credentials)
	if err != nil {
		a.logger.Error("", zap.Error(err))
	}

	output, err := a.service.Authenticate(r.Context(), credentials)
	if err != nil {
		a.logger.Error("authenticate failed", zap.Error(err))
	}

	var resp = struct {
		AccessToken  *string `json:"accessToken"`
		ExpiresIn    int32   `json:"expiresIn"`
		IdToken      *string `json:"idToken"`
		RefreshToken *string `json:"refreshToken"`
		TokenType    *string `json:"tokenType"`
	}{
		AccessToken:  output.AuthenticationResult.AccessToken,
		ExpiresIn:    output.AuthenticationResult.ExpiresIn,
		IdToken:      output.AuthenticationResult.IdToken,
		RefreshToken: output.AuthenticationResult.RefreshToken,
		TokenType:    output.AuthenticationResult.TokenType,
	}

	err = web.Respond(r.Context(), w, resp, http.StatusOK)
	if err != nil {
		a.logger.Error("", zap.Error(err))
	}
}

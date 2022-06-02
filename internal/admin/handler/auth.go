package handler

import (
	"context"
	"net/http"

	"github.com/devpies/core/pkg/web"

	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"go.uber.org/zap"
)

type authService interface {
	Authenticate(ctx context.Context, email, password string) (*cip.AdminInitiateAuthOutput, error)
	GetSubject(ctx context.Context, token []byte) (string, error)
	RespondToNewPasswordRequiredChallenge(ctx context.Context, email, password string, session string) (*cip.AdminRespondToAuthChallengeOutput, error)
}

type AuthHandler struct {
	logger  *zap.Logger
	service authService
}

func NewAuthHandler(logger *zap.Logger, service authService) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: service,
	}
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// AuthenticateCredentials handles email and password values from the admin login form.
func (ah *AuthHandler) AuthenticateCredentials(w http.ResponseWriter, r *http.Request) error {
	var err error

	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = web.Decode(r, &payload)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	output, err := ah.service.Authenticate(r.Context(), payload.Email, payload.Password)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	if output.AuthenticationResult != nil {
		var resp = struct {
			IDToken *string `json:"idToken"`
		}{
			IDToken: output.AuthenticationResult.IdToken,
		}
		return web.Respond(r.Context(), w, resp, http.StatusOK)
	}

	var resp = struct {
		ChallengeName types.ChallengeNameType `json:"challengeName"`
		Session       *string                 `json:"session"`
	}{
		ChallengeName: output.ChallengeName,
		Session:       output.Session,
	}

	return web.Respond(r.Context(), w, resp, http.StatusOK)
}

func (ah *AuthHandler) SetupNewUserWithSecurePassword(w http.ResponseWriter, r *http.Request) error {
	var err error

	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Session  string `json:"session"`
	}

	err = web.Decode(r, &payload)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	output, err := ah.service.RespondToNewPasswordRequiredChallenge(r.Context(), payload.Email, payload.Password, payload.Session)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	if output.AuthenticationResult.IdToken == nil {
		return web.NewRequestError(err, http.StatusInternalServerError)
	}
	token := *output.AuthenticationResult.IdToken

	sub, err := ah.service.GetSubject(r.Context(), []byte(token))
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	var resp = struct {
		IDToken string `json:"idToken"`
		UserID  string `json:"userID"`
	}{
		IDToken: token,
		UserID:  sub,
	}

	return web.Respond(r.Context(), w, resp, http.StatusOK)
}

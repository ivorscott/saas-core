// Package handler manages the presentation layer for handling incoming requests.
package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"

	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"go.uber.org/zap"
)

type authService interface {
	Authenticate(ctx context.Context, email, password string) (*cip.AdminInitiateAuthOutput, error)
	CreateUserSession(ctx context.Context, token []byte) error
	CreatePasswordChallengeSession(ctx context.Context)
	RespondToNewPasswordRequiredChallenge(ctx context.Context, email, password string, session string) (*cip.AdminRespondToAuthChallengeOutput, error)
}

type sessionManager interface {
	Destroy(ctx context.Context) error
	RenewToken(ctx context.Context) error
}

// AuthHandler contains various auth related handlers.
type AuthHandler struct {
	logger      *zap.Logger
	render      renderer
	session     sessionManager
	authService authService
}

var (
	// ErrIncorrectUsernameOrPassword represents an AWS Cognito error caused by invalid credentials.
	ErrIncorrectUsernameOrPassword = errors.New("incorrect username or password")
	// ErrPasswordAttemptsExceeded represents an AWS Cognito error caused by exceeding allowed password attempts.
	ErrPasswordAttemptsExceeded = errors.New("password attempts exceeded")
	// ErrNotAuthorizedException represents an unknown AWS Cognito error effecting login.
	ErrNotAuthorizedException = errors.New("login failed")
)

// NewAuthHandler returns a new authentication handler.
func NewAuthHandler(logger *zap.Logger, renderEngine renderer, session sessionManager, authService authService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		render:      renderEngine,
		session:     session,
		authService: authService,
	}
}

// VerifyTokenNoop does nothing. It's only used to allow the app middleware to verify the id_token on each page.
func (ah *AuthHandler) VerifyTokenNoop(w http.ResponseWriter, r *http.Request) error {
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// LoginPage displays a form to allow users to sign in.
func (ah *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) error {
	return ah.render.Template(w, r, "login", nil)
}

// ForceNewPasswordPage displays a form where freshly onboarded users can change their OTP.
func (ah *AuthHandler) ForceNewPasswordPage(w http.ResponseWriter, r *http.Request) error {
	return ah.render.Template(w, r, "force-new-password", nil)
}

// Logout allows users to log out by destroying the existing session.
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	var err error

	err = ah.session.Destroy(r.Context())
	if err != nil {
		ah.logger.Error("failure on session destroy", zap.Error(err))
		return web.NewShutdownError(err.Error())
	}

	// Renew the session token everytime a user logs out.
	err = ah.session.RenewToken(r.Context())
	if err != nil {
		ah.logger.Error("failure on session renewal", zap.Error(err))
		return web.NewShutdownError(err.Error())
	}

	web.Redirect(w, r, "/", http.StatusSeeOther)
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// AuthenticateCredentials handles email and password values from the admin login form.
func (ah *AuthHandler) AuthenticateCredentials(w http.ResponseWriter, r *http.Request) error {
	var (
		err     error
		payload model.AuthCredentials
	)

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	// Renew the session token everytime a user logs in.
	err = ah.session.RenewToken(r.Context())
	if err != nil {
		ah.logger.Error("failure on session renewal", zap.Error(err))
		return web.NewShutdownError(err.Error())
	}

	// Authenticate.
	output, err := ah.authService.Authenticate(r.Context(), payload.Email, payload.Password)
	if err != nil {
		ah.logger.Info("", zap.Error(err))
		switch {
		case strings.Contains(strings.ToLower(err.Error()), ErrIncorrectUsernameOrPassword.Error()):
			err = ErrIncorrectUsernameOrPassword
		case strings.Contains(strings.ToLower(err.Error()), ErrPasswordAttemptsExceeded.Error()):
			err = ErrPasswordAttemptsExceeded
		default:
			err = ErrNotAuthorizedException
		}
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	// On success.
	if output.AuthenticationResult != nil {
		err = ah.authService.CreateUserSession(r.Context(), []byte(*output.AuthenticationResult.IdToken))
		if err != nil {
			ah.logger.Error("failure on user session creation", zap.Error(err))
			return web.NewShutdownError(err.Error())
		}

		var resp = struct {
			IDToken string `json:"idToken"`
		}{
			IDToken: *output.AuthenticationResult.IdToken,
		}

		return web.Respond(r.Context(), w, resp, http.StatusOK)
	}

	ah.authService.CreatePasswordChallengeSession(r.Context())

	// On challenge.
	var resp = struct {
		ChallengeName types.ChallengeNameType `json:"challengeName"`
		Session       string                  `json:"session"`
	}{
		ChallengeName: output.ChallengeName,
		Session:       *output.Session,
	}

	return web.Respond(r.Context(), w, resp, http.StatusOK)
}

// SetupNewUserWithSecurePassword responds to force change password challenge.
func (ah *AuthHandler) SetupNewUserWithSecurePassword(w http.ResponseWriter, r *http.Request) error {
	var (
		err     error
		payload struct {
			model.AuthCredentials
			Session string `json:"session" validate:"required"`
		}
	)

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	output, err := ah.authService.RespondToNewPasswordRequiredChallenge(r.Context(), payload.Email, payload.Password, payload.Session)
	if err != nil {
		ah.logger.Info("failure on password required challenge response", zap.Error(err))
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	err = ah.authService.CreateUserSession(r.Context(), []byte(*output.AuthenticationResult.IdToken))
	if err != nil {
		ah.logger.Error("failure on user session creation", zap.Error(err))
		return web.NewShutdownError(err.Error())
	}

	var resp = struct {
		IDToken string `json:"idToken"`
	}{
		IDToken: *output.AuthenticationResult.IdToken,
	}

	return web.Respond(r.Context(), w, resp, http.StatusOK)
}

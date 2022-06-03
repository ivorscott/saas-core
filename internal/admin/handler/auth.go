package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/devpies/core/internal/admin/config"
	"github.com/devpies/core/internal/admin/model"
	"github.com/devpies/core/internal/admin/render"
	"github.com/devpies/core/pkg/web"

	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"go.uber.org/zap"
)

// AuthHandler contains various auth related handlers.
type AuthHandler struct {
	logger  *zap.Logger
	config  config.Config
	render  *render.Render
	service authService
	session *scs.SessionManager
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
func NewAuthHandler(logger *zap.Logger, config config.Config, renderEngine *render.Render, service authService, session *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		config:  config,
		render:  renderEngine,
		service: service,
		session: session,
	}
}

// Login displays a form to allow users to sign in.
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	return ah.render.Template(w, r, "login", nil)
}

// ForceNewPassword displays a form where freshly onboarded users can change their OTP.
func (ah *AuthHandler) ForceNewPassword(w http.ResponseWriter, r *http.Request) error {
	return ah.render.Template(w, r, "new-password", nil)
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

	web.SetContextStatusCode(r.Context(), http.StatusOK)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

// AuthenticateCredentials handles email and password values from the admin login form.
func (ah *AuthHandler) AuthenticateCredentials(w http.ResponseWriter, r *http.Request) error {
	var (
		err     error
		payload model.AuthCredentials
	)

	// Renew the session token everytime a user logs in.
	err = ah.session.RenewToken(r.Context())
	if err != nil {
		ah.logger.Error("failure on session renewal", zap.Error(err))
		return web.NewShutdownError(err.Error())
	}

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	// Authenticate.
	output, err := ah.service.Authenticate(r.Context(), payload.Email, payload.Password)
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
		err = ah.service.CreateUserSession(r.Context(), []byte(*output.AuthenticationResult.IdToken))
		if err != nil {
			ah.logger.Error("failure on user session creation", zap.Error(err))
			return web.NewShutdownError(err.Error())
		}

		var resp = struct {
			IDToken *string `json:"idToken"`
		}{
			IDToken: output.AuthenticationResult.IdToken,
		}

		return web.Respond(r.Context(), w, resp, http.StatusOK)
	}

	// On challenge.
	var resp = struct {
		ChallengeName types.ChallengeNameType `json:"challengeName"`
		Session       *string                 `json:"session"`
	}{
		ChallengeName: output.ChallengeName,
		Session:       output.Session,
	}

	return web.Respond(r.Context(), w, resp, http.StatusOK)
}

// SetupNewUserWithSecurePassword responds to force change password challenge.
func (ah *AuthHandler) SetupNewUserWithSecurePassword(w http.ResponseWriter, r *http.Request) error {
	var (
		err     error
		payload struct {
			model.AuthCredentials
			Session string `json:"session"`
		}
	)

	err = web.Decode(r, &payload)
	if err != nil {
		return err
	}

	output, err := ah.service.RespondToNewPasswordRequiredChallenge(r.Context(), payload.Email, payload.Password, payload.Session)
	if err != nil {
		ah.logger.Info("failure on password required challenge response", zap.Error(err))
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	err = ah.service.CreateUserSession(r.Context(), []byte(*output.AuthenticationResult.IdToken))
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

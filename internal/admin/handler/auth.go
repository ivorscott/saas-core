package handler

import (
	"net/http"

	"github.com/devpies/core/internal/admin/config"
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
func (app *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := app.render.Template(w, r, "login", nil); err != nil {
		app.logger.Error("login", zap.Error(err))
	}
}

// ForceNewPassword displays a form where freshly onboarded users can change their OTP.
func (app *AuthHandler) ForceNewPassword(w http.ResponseWriter, r *http.Request) {
	if err := app.render.Template(w, r, "new-password", nil); err != nil {
		app.logger.Error("new-password", zap.Error(err))
	}
}

// Logout allows users to log out by destroying the existing session.
func (app *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var err error

	err = app.session.Destroy(r.Context())
	if err != nil {
		app.logger.Error("session destroy failed", zap.Error(err))
	}

	// Renew the session token everytime a user logs out.
	err = app.session.RenewToken(r.Context())
	if err != nil {
		app.logger.Error("session renew token failed", zap.Error(err))
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AuthenticateCredentials handles email and password values from the admin login form.
func (app *AuthHandler) AuthenticateCredentials(w http.ResponseWriter, r *http.Request) {
	var err error

	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Renew the session token everytime a user logs in.
	err = app.session.RenewToken(r.Context())
	if err != nil {
		app.logger.Error("error on renew session token", zap.Error(err))
		_ = web.Respond(r.Context(), w, nil, http.StatusInternalServerError)
		return
	}
	err = web.Decode(r, &payload)
	if err != nil {
		_ = web.Respond(r.Context(), w, nil, http.StatusBadRequest)
		return
	}
	output, err := app.service.Authenticate(r.Context(), payload.Email, payload.Password)
	if err != nil {
		_ = web.Respond(r.Context(), w, nil, http.StatusUnauthorized)
		return
	}
	if output.AuthenticationResult != nil {
		var resp = struct {
			IDToken *string `json:"idToken"`
		}{
			IDToken: output.AuthenticationResult.IdToken,
		}

		err = app.service.CreateUserSession(r.Context(), []byte(*output.AuthenticationResult.IdToken))
		if err != nil {
			app.logger.Error("error creating user session")
			_ = web.Respond(r.Context(), w, nil, http.StatusInternalServerError)
			return
		}
		_ = web.Respond(r.Context(), w, resp, http.StatusOK)
		return
	}

	var resp = struct {
		ChallengeName types.ChallengeNameType `json:"challengeName"`
		Session       *string                 `json:"session"`
	}{
		ChallengeName: output.ChallengeName,
		Session:       output.Session,
	}

	_ = web.Respond(r.Context(), w, resp, http.StatusOK)
}

// SetupNewUserWithSecurePassword responds to force change password challenge.
func (app *AuthHandler) SetupNewUserWithSecurePassword(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Session  string `json:"session"`
		}
	)

	err = web.Decode(r, &payload)
	if err != nil {
		app.logger.Error("error decoding payload", zap.Error(err))
		_ = web.Respond(r.Context(), w, nil, http.StatusBadRequest)
		return
	}

	output, err := app.service.RespondToNewPasswordRequiredChallenge(r.Context(), payload.Email, payload.Password, payload.Session)
	if err != nil {
		app.logger.Error("error responding to challenge", zap.Error(err))
		_ = web.Respond(r.Context(), w, nil, http.StatusInternalServerError)
		return
	}

	err = app.service.CreateUserSession(r.Context(), []byte(*output.AuthenticationResult.IdToken))
	if err != nil {
		app.logger.Error("error creating user session")
		_ = web.Respond(r.Context(), w, nil, http.StatusInternalServerError)
		return
	}

	var resp = struct {
		IDToken string `json:"idToken"`
	}{
		IDToken: *output.AuthenticationResult.IdToken,
	}

	_ = web.Respond(r.Context(), w, resp, http.StatusOK)
}

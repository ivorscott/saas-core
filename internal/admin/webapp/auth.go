package webapp

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/devpies/core/pkg/web"

	"go.uber.org/zap"
)

func (app *WebApp) Login(w http.ResponseWriter, r *http.Request) {
	if err := app.render.Template(w, r, "login", nil); err != nil {
		app.logger.Error("login", zap.Error(err))
	}
}

func (app *WebApp) Logout(w http.ResponseWriter, r *http.Request) {
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

func (app *WebApp) ForceNewPassword(w http.ResponseWriter, r *http.Request) {
	if err := app.render.Template(w, r, "new-password", nil); err != nil {
		app.logger.Error("new-password", zap.Error(err))
	}
}

// AuthenticateCredentials handles email and password values from the admin login form.
func (app *WebApp) AuthenticateCredentials(w http.ResponseWriter, r *http.Request) {
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

func (app *WebApp) SetupNewUserWithSecurePassword(w http.ResponseWriter, r *http.Request) {
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

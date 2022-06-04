package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/devpies/saas-core/internal/admin/config"
	"github.com/devpies/saas-core/internal/admin/handler"
	"github.com/devpies/saas-core/internal/admin/mocks"
	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthHandler_LoginPage(t *testing.T) {
	basePath := "/"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupAuthHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.render.On("Template", mock.Anything, mock.AnythingOfType("*http.Request"), "login", testTemplateData).Return(nil)

		handle.ServeHTTP(w, r)

		deps.render.AssertExpectations(t)
	})
}

func TestAuthHandler_ForceNewPasswordPage(t *testing.T) {
	basePath := "/force-new-password"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupAuthHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.render.On("Template", mock.Anything, mock.AnythingOfType("*http.Request"), "force-new-password", testTemplateData).Return(nil)

		handle.ServeHTTP(w, r)

		deps.render.AssertExpectations(t)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	basePath := "/admin/logout"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupAuthHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.session.On("Destroy", mockCtx).Return(nil)
		deps.session.On("RenewToken", mockCtx).Return(nil)

		handle.ServeHTTP(w, r)

		deps.session.AssertExpectations(t)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	})

	t.Run("500 on logout", func(t *testing.T) {
		tests := []struct {
			name         string
			expectations func(deps authHandlerDeps)
		}{
			{
				name: "error on session destroy failure",
				expectations: func(deps authHandlerDeps) {
					deps.session.On("Destroy", mockCtx).Return(nil)
					deps.session.On("RenewToken", mockCtx).Return(assert.AnError)
				},
			},
			{
				name: "error on session renew token failure",
				expectations: func(deps authHandlerDeps) {
					deps.session.On("Destroy", mockCtx).Return(assert.AnError)
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				handle, deps := setupAuthHandler()

				r := httptest.NewRequest(http.MethodGet, basePath, nil)
				w := httptest.NewRecorder()

				tc.expectations(deps)

				handle.ServeHTTP(w, r)

				var resp web.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)

				deps.session.AssertExpectations(t)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, http.StatusText(http.StatusInternalServerError), resp.Error)
			})
		}
	})
}

func TestAuthHandler_AuthenticateCredentials(t *testing.T) {
	basePath := "/authenticate"

	t.Run("200 on authenticated", func(t *testing.T) {
		var (
			adminInitiateAuthOutput = &cognitoidentityprovider.AdminInitiateAuthOutput{
				AuthenticationResult: &types.AuthenticationResultType{
					IdToken: &testIDToken,
				},
			}
			nilAuthenticationResult = &cognitoidentityprovider.AdminInitiateAuthOutput{
				AuthenticationResult: nil,
				ChallengeName:        "NEW_PASSWORD_REQUIRED",
				Session:              &testPChalSession,
			}
			request = model.AuthCredentials{
				Email:    testEmail,
				Password: testPassword,
			}
		)

		type loginResponse struct {
			IDToken string `json:"idToken"`
		}
		type pChalResponse struct {
			ChallengeName types.ChallengeNameType `json:"challengeName"`
			Session       string                  `json:"session"`
		}

		tests := []struct {
			name         string
			expectations func(deps authHandlerDeps)
			respJSON     func() []byte
		}{
			{
				name: "login",
				expectations: func(deps authHandlerDeps) {
					deps.session.On("RenewToken", mockCtx).Return(nil)
					deps.authService.On("Authenticate", mockCtx, testEmail, testPassword).Return(adminInitiateAuthOutput, nil)
					deps.authService.On("CreateUserSession", mockCtx, []byte(testIDToken)).Return(nil)
				},
				respJSON: func() []byte {
					b, _ := json.Marshal(loginResponse{testIDToken})
					return b
				},
			},
			{
				name: "password challenge",
				expectations: func(deps authHandlerDeps) {
					deps.session.On("RenewToken", mockCtx).Return(nil)
					deps.authService.On("Authenticate", mockCtx, testEmail, testPassword).Return(nilAuthenticationResult, nil)
					deps.authService.On("CreatePasswordChallengeSession", mockCtx)
				},
				respJSON: func() []byte {
					b, _ := json.Marshal(pChalResponse{"NEW_PASSWORD_REQUIRED", testPChalSession})
					return b
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				handle, deps := setupAuthHandler()

				payload, err := json.Marshal(request)
				assert.Nil(t, err)

				r := httptest.NewRequest(http.MethodGet, basePath, bytes.NewReader(payload))
				w := httptest.NewRecorder()

				tc.expectations(deps)

				handle.ServeHTTP(w, r)

				deps.session.AssertExpectations(t)
				deps.authService.AssertExpectations(t)

				assert.True(t, assert.Exactlyf(t, tc.respJSON(), w.Body.Bytes(), "not exact match"))
				assert.Equal(t, http.StatusOK, w.Code)
			})
		}
	})

	t.Run("400 on invalid payload", func(t *testing.T) {
		var (
			request = model.AuthCredentials{
				Email:    "",
				Password: "",
			}
			resp web.ErrorResponse
			err  error
		)

		handle, _ := setupAuthHandler()

		payload, err := json.Marshal(request)
		assert.Nil(t, err)

		r := httptest.NewRequest(http.MethodGet, basePath, bytes.NewReader(payload))
		w := httptest.NewRecorder()

		handle.ServeHTTP(w, r)

		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, fieldValidationErr, resp.Error)
	})

	t.Run("401 on unauthorized login", func(t *testing.T) {
		var (
			validRequest = model.AuthCredentials{
				Email:    testEmail,
				Password: testPassword,
			}
			resp web.ErrorResponse
		)

		tests := []struct {
			name        string
			errToReturn error
			errResp     string
		}{
			{
				name:        "incorrect username and password error",
				errToReturn: handler.ErrIncorrectUsernameOrPassword,
				errResp:     handler.ErrIncorrectUsernameOrPassword.Error(),
			},
			{
				name:        "password attempts exceeded error",
				errToReturn: handler.ErrPasswordAttemptsExceeded,
				errResp:     handler.ErrPasswordAttemptsExceeded.Error(),
			},
			{
				name:        "generic unauthorized error",
				errToReturn: assert.AnError,
				errResp:     handler.ErrNotAuthorizedException.Error(),
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				handle, deps := setupAuthHandler()

				payload, err := json.Marshal(validRequest)
				assert.Nil(t, err)

				r := httptest.NewRequest(http.MethodGet, basePath, bytes.NewReader(payload))
				w := httptest.NewRecorder()

				deps.session.
					On("RenewToken", mock.AnythingOfType("*context.valueCtx")).
					Return(nil)
				deps.authService.
					On("Authenticate", mockCtx, testEmail, testPassword).
					Return(nil, tc.errToReturn)

				handle.ServeHTTP(w, r)

				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)

				deps.session.AssertExpectations(t)

				assert.Equal(t, http.StatusUnauthorized, w.Code)
				assert.Equal(t, tc.errResp, resp.Error)
			})
		}
	})

	t.Run("500 on authenticate", func(t *testing.T) {
		var (
			adminInitiateAuthOutput = &cognitoidentityprovider.AdminInitiateAuthOutput{
				AuthenticationResult: &types.AuthenticationResultType{
					IdToken: &testIDToken,
				},
			}
			request = model.AuthCredentials{
				Email:    testEmail,
				Password: testPassword,
			}
			resp web.ErrorResponse
		)

		tests := []struct {
			name         string
			expectations func(deps authHandlerDeps)
		}{
			{
				name: "renew session token failure",
				expectations: func(deps authHandlerDeps) {
					deps.session.On("RenewToken", mock.AnythingOfType("*context.valueCtx")).Return(assert.AnError)
				},
			},
			{
				name: "create user session failure",
				expectations: func(deps authHandlerDeps) {
					deps.session.On("RenewToken", mockCtx).Return(nil)
					deps.authService.On("Authenticate", mockCtx, testEmail, testPassword).Return(adminInitiateAuthOutput, nil)
					deps.authService.On("CreateUserSession", mockCtx, []byte(testIDToken)).Return(assert.AnError)
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				handle, deps := setupAuthHandler()

				payload, err := json.Marshal(request)
				assert.Nil(t, err)

				r := httptest.NewRequest(http.MethodGet, basePath, bytes.NewReader(payload))
				w := httptest.NewRecorder()

				tc.expectations(deps)

				handle.ServeHTTP(w, r)

				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)

				deps.session.AssertExpectations(t)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, http.StatusText(http.StatusInternalServerError), resp.Error)
			})
		}
	})
}

type authHandlerDeps struct {
	logger      *zap.Logger
	render      *mocks.Renderer
	session     *mocks.SessionManager
	authService *mocks.AuthService
}

func setupAuthHandler() (http.Handler, authHandlerDeps) {
	router := chi.NewRouter()
	shutdown := make(chan os.Signal, 1)
	logger := zap.NewNop()
	renderEngine := &mocks.Renderer{}
	authService := &mocks.AuthService{}
	sessionManager := &mocks.SessionManager{}

	auth := handler.NewAuthHandler(logger, config.Config{}, renderEngine, sessionManager, authService)

	app := web.NewApp(router, shutdown, logger, []web.Middleware{mid.Errors(logger)}...)

	app.Handle(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		return auth.LoginPage(w, r)
	})

	app.Handle(http.MethodGet, "/force-new-password", func(w http.ResponseWriter, r *http.Request) error {
		return auth.ForceNewPasswordPage(w, r)
	})

	app.Handle(http.MethodGet, "/authenticate", func(w http.ResponseWriter, r *http.Request) error {
		return auth.AuthenticateCredentials(w, r)
	})

	app.Handle(http.MethodGet, "/admin/logout", func(w http.ResponseWriter, r *http.Request) error {
		return auth.Logout(w, r)
	})

	return router, authHandlerDeps{logger, renderEngine, sessionManager, authService}
}

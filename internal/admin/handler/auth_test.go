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

func TestAuthHandler_Login(t *testing.T) {
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

func TestAuthHandler_ForceNewPassword(t *testing.T) {
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

	t.Run("500 error on session destroy failure", func(t *testing.T) {
		handle, deps := setupAuthHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.session.On("Destroy", mockCtx).Return(nil)
		deps.session.On("RenewToken", mockCtx).Return(assert.AnError)

		handle.ServeHTTP(w, r)

		var resp web.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)

		deps.session.AssertExpectations(t)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), resp.Error)
	})

	t.Run("500 error on session renew token failure", func(t *testing.T) {
		handle, deps := setupAuthHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.session.On("Destroy", mockCtx).Return(assert.AnError)

		handle.ServeHTTP(w, r)

		var resp web.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)

		deps.session.AssertExpectations(t)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), resp.Error)
	})
}

func TestAuthHandler_AuthenticateCredentials(t *testing.T) {
	basePath := "/authenticate"

	t.Run("success", func(t *testing.T) {
		var err error

		handle, deps := setupAuthHandler()

		payload := model.AuthCredentials{
			Email:    testEmail,
			Password: testPassword,
		}
		data, err := json.Marshal(payload)
		assert.Nil(t, err)

		r := httptest.NewRequest(http.MethodGet, basePath, bytes.NewReader(data))
		w := httptest.NewRecorder()

		deps.session.On("RenewToken", mockCtx).Return(nil)
		deps.authService.On("Authenticate", mockCtx, testEmail, testPassword).
			Return(&cognitoidentityprovider.AdminInitiateAuthOutput{
				AuthenticationResult: &types.AuthenticationResultType{
					IdToken: &testIDToken,
				},
			}, nil)
		deps.authService.On("CreateUserSession", mockCtx, []byte(testIDToken)).Return(nil)
		handle.ServeHTTP(w, r)

		deps.session.AssertExpectations(t)

		var resp struct {
			IDToken string `json:"idToken"`
		}

		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)

		deps.session.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, testIDToken, resp.IDToken)
	})

	t.Run("500 error on session renew token failure", func(t *testing.T) {
		handle, deps := setupAuthHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.session.On("RenewToken", mock.AnythingOfType("*context.valueCtx")).Return(assert.AnError)

		handle.ServeHTTP(w, r)

		var resp web.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)

		deps.session.AssertExpectations(t)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), resp.Error)
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
		return auth.Login(w, r)
	})

	app.Handle(http.MethodGet, "/force-new-password", func(w http.ResponseWriter, r *http.Request) error {
		return auth.ForceNewPassword(w, r)
	})

	app.Handle(http.MethodGet, "/authenticate", func(w http.ResponseWriter, r *http.Request) error {
		return auth.AuthenticateCredentials(w, r)
	})

	app.Handle(http.MethodGet, "/admin/logout", func(w http.ResponseWriter, r *http.Request) error {
		return auth.Logout(w, r)
	})

	return router, authHandlerDeps{logger, renderEngine, sessionManager, authService}
}

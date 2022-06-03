package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devpies/core/internal/admin/handler"
	"github.com/devpies/core/internal/admin/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"

	"go.uber.org/zap"
)

func TestWebPageHandler_Dashboard(t *testing.T) {
	basePath := "/admin"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupWebPageHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.render.On("Template", mock.Anything, mock.AnythingOfType("*http.Request"), "dashboard", testTemplateData).Return(nil)

		handle.ServeHTTP(w, r)

		deps.render.AssertExpectations(t)
	})
}

func TestWebPageHandler_Tenants(t *testing.T) {
	basePath := "/admin/tenants"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupWebPageHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.render.On("Template", mock.Anything, mock.AnythingOfType("*http.Request"), "tenants", testTemplateData).Return(nil)

		handle.ServeHTTP(w, r)

		deps.render.AssertExpectations(t)
	})
}

func TestWebPageHandler_CreateTenant(t *testing.T) {
	basePath := "/admin/create-tenant"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupWebPageHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.render.On("Template", mock.Anything, mock.AnythingOfType("*http.Request"), "create-tenant", testTemplateData).Return(nil)

		handle.ServeHTTP(w, r)

		deps.render.AssertExpectations(t)
	})
}

func TestWebPageHandler_E404(t *testing.T) {
	basePath := "/noSuchThing"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupWebPageHandler()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.render.On("Template", mock.Anything, mock.AnythingOfType("*http.Request"), "404", testTemplateData).Return(nil)
		deps.setStatus.On("Execute", mock.AnythingOfType("*context.valueCtx"), http.StatusNotFound)

		handle.ServeHTTP(w, r)

		deps.render.AssertExpectations(t)
	})
}

type webPageHandlerDeps struct {
	logger    *zap.Logger
	render    *mocks.Renderer
	setStatus *mocks.StatusSetter
}

func setupWebPageHandler() (http.Handler, webPageHandlerDeps) {
	router := chi.NewRouter()
	logger := zap.NewNop()
	renderEngine := &mocks.Renderer{}
	setStatus := &mocks.StatusSetter{}
	webpage := handler.NewWebPageHandler(logger, renderEngine, setStatus.Execute)

	router.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		_ = webpage.Dashboard(w, r)
	})

	router.Get("/admin/tenants", func(w http.ResponseWriter, r *http.Request) {
		_ = webpage.Tenants(w, r)
	})

	router.Get("/admin/create-tenant", func(w http.ResponseWriter, r *http.Request) {
		_ = webpage.CreateTenant(w, r)
	})

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		_ = webpage.E404(w, r)
	})

	return router, webPageHandlerDeps{logger, renderEngine, setStatus}
}

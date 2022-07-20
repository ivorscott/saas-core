package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/devpies/saas-core/internal/project/handler"
	"github.com/devpies/saas-core/internal/project/mocks"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/internal/project/res/testutils"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestProjectHandler_Create(t *testing.T) {
	basePath := "/projects"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupProjectRouter()

		np := model.NewProject{
			Name: "My Project",
		}
		project := model.Project{
			ID:   testutils.MockUUID,
			Name: "My Project",
		}

		b, err := json.Marshal(&np)
		assert.Nil(t, err)

		r := httptest.NewRequest(http.MethodPost, basePath, bytes.NewReader(b))
		w := httptest.NewRecorder()

		deps.projectService.On("Create", mock.AnythingOfType("*context.valueCtx"), np, mock.AnythingOfType("time.Time")).Return(project, nil)
		deps.columnService.On("CreateColumns", mock.AnythingOfType("*context.valueCtx"), project.ID, mock.AnythingOfType("time.Time")).Return(nil)

		handle.ServeHTTP(w, r)

		expectedProject, err := json.Marshal(&project)
		assert.Nil(t, err)
		assert.Equal(t, expectedProject, w.Body.Bytes())
		assert.Equal(t, http.StatusCreated, w.Code)
		deps.projectService.AssertExpectations(t)
		deps.columnService.AssertExpectations(t)
	})

	t.Run("error 400", func(t *testing.T) {
		handle, _ := setupProjectRouter()

		np := model.NewProject{}
		response := web.ErrorResponse{
			Error: "field validation error",
		}

		b, err := json.Marshal(&np)
		assert.Nil(t, err)

		r := httptest.NewRequest(http.MethodPost, basePath, bytes.NewReader(b))
		w := httptest.NewRecorder()

		handle.ServeHTTP(w, r)

		expectedResponse, err := json.Marshal(&response)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, w.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error 500 project service", func(t *testing.T) {
		handle, deps := setupProjectRouter()

		np := model.NewProject{
			Name: "My Project",
		}
		response := web.ErrorResponse{
			Error: http.StatusText(http.StatusInternalServerError),
		}

		b, err := json.Marshal(&np)
		assert.Nil(t, err)

		r := httptest.NewRequest(http.MethodPost, basePath, bytes.NewReader(b))
		w := httptest.NewRecorder()

		deps.projectService.
			On("Create", mock.AnythingOfType("*context.valueCtx"), np, mock.AnythingOfType("time.Time")).
			Return(model.Project{}, assert.AnError)

		handle.ServeHTTP(w, r)

		expectedResponse, err := json.Marshal(&response)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, w.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		deps.projectService.AssertExpectations(t)
	})

	t.Run("error 500 column service", func(t *testing.T) {
		handle, deps := setupProjectRouter()

		np := model.NewProject{
			Name: "My Project",
		}
		project := model.Project{
			ID:   testutils.MockUUID,
			Name: "My Project",
		}
		response := web.ErrorResponse{
			Error: http.StatusText(http.StatusInternalServerError),
		}

		b, err := json.Marshal(&np)
		assert.Nil(t, err)

		r := httptest.NewRequest(http.MethodPost, basePath, bytes.NewReader(b))
		w := httptest.NewRecorder()

		deps.projectService.
			On("Create", mock.AnythingOfType("*context.valueCtx"), np, mock.AnythingOfType("time.Time")).
			Return(project, nil)

		deps.columnService.
			On("CreateColumns", mock.AnythingOfType("*context.valueCtx"), project.ID, mock.AnythingOfType("time.Time")).
			Return(assert.AnError)

		handle.ServeHTTP(w, r)

		expectedResponse, err := json.Marshal(&response)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, w.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		deps.projectService.AssertExpectations(t)
	})
}

func TestProjectHandler_List(t *testing.T) {
	basePath := "/projects"

	t.Run("success", func(t *testing.T) {
		var projects []model.Project

		handle, deps := setupProjectRouter()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		w := httptest.NewRecorder()

		deps.projectService.On("List", mock.AnythingOfType("*context.valueCtx"), false).Return(projects, nil)

		handle.ServeHTTP(w, r)

		expectedProjects, err := json.Marshal(&projects)
		assert.Nil(t, err)
		assert.Equal(t, expectedProjects, w.Body.Bytes())
		assert.Equal(t, http.StatusOK, w.Code)
		deps.projectService.AssertExpectations(t)
	})

	t.Run("success all tenants", func(t *testing.T) {
		var projects []model.Project

		handle, deps := setupProjectRouter()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		r.Header.Set("BasePath", "projects")

		w := httptest.NewRecorder()

		deps.projectService.On("List", mock.AnythingOfType("*context.valueCtx"), true).Return(projects, nil)

		handle.ServeHTTP(w, r)

		expectedProjects, err := json.Marshal(&projects)
		assert.Nil(t, err)
		assert.Equal(t, w.Body.Bytes(), expectedProjects)
		assert.Equal(t, w.Code, http.StatusOK)
		deps.projectService.AssertExpectations(t)
	})

	t.Run("error 500", func(t *testing.T) {
		response := web.ErrorResponse{
			Error: http.StatusText(http.StatusInternalServerError),
		}

		handle, deps := setupProjectRouter()

		r := httptest.NewRequest(http.MethodGet, basePath, nil)
		r.Header.Set("BasePath", "projects")

		w := httptest.NewRecorder()

		deps.projectService.On("List", mock.AnythingOfType("*context.valueCtx"), true).Return(nil, assert.AnError)

		handle.ServeHTTP(w, r)

		expectedResponse, err := json.Marshal(&response)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, w.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		deps.projectService.AssertExpectations(t)
	})
}

type projectHandlerDeps struct {
	logger         *zap.Logger
	projectService *mocks.ProjectService
	columnService  *mocks.ColumnService
	taskService    *mocks.TaskService
}

func setupProjectRouter() (http.Handler, projectHandlerDeps) {
	router := chi.NewRouter()
	logger := zap.NewNop()
	projectService := &mocks.ProjectService{}
	columnService := &mocks.ColumnService{}
	taskService := &mocks.TaskService{}
	shutdown := make(chan os.Signal, 1)

	middleware := []web.Middleware{
		mid.Logger(logger),
		mid.Errors(logger),
		mid.Panics(logger),
	}

	projects := handler.NewProjectHandler(logger, projectService, columnService, taskService)

	app := web.NewApp(router, shutdown, logger, middleware...)
	app.Handle(http.MethodPost, "/projects", projects.Create)
	app.Handle(http.MethodGet, "/projects", projects.List)

	return router, projectHandlerDeps{logger, projectService, columnService, taskService}
}

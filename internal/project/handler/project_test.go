package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devpies/saas-core/internal/project/handler"
	"github.com/devpies/saas-core/internal/project/mocks"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/internal/project/res/testutils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectHandler_Create(t *testing.T) {
	basePath := "/projects"

	t.Run("success", func(t *testing.T) {
		handle, deps := setupProjectRouter()

		np := model.NewProject{
			Name: "My Project",
		}
		expectedProject := model.Project{
			ID:   testutils.MockUUID,
			Name: "My Project",
		}

		b, _ := json.Marshal(&np)

		r := httptest.NewRequest(http.MethodPost, basePath, bytes.NewReader(b))
		w := httptest.NewRecorder()

		deps.projectService.On("Create", mock.AnythingOfType("*context.valueCtx"), np, mock.AnythingOfType("time.Time")).Return(expectedProject, nil)

		for i, title := range [4]string{"To Do", "In Progress", "Review", "Done"} {
			nt := model.NewColumn{
				ProjectID:  expectedProject.ID,
				Title:      title,
				ColumnName: fmt.Sprintf(`column-%d`, i+1),
			}
			c := model.Column{
				ProjectID:  expectedProject.ID,
				Title:      title,
				ColumnName: fmt.Sprintf(`column-%d`, i+1),
			}
			deps.columnService.On("Create", mock.AnythingOfType("*context.valueCtx"), nt, mock.AnythingOfType("time.Time")).Return(c, nil)
		}

		handle.ServeHTTP(w, r)

		deps.projectService.AssertExpectations(t)
		deps.columnService.AssertExpectations(t)
		assert.Equal(t, w.Code, http.StatusCreated)
		project, _ := json.Marshal(&expectedProject)
		assert.Equal(t, w.Body.Bytes(), project)
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

	projects := handler.NewProjectHandler(logger, projectService, columnService, taskService)

	router.Post("/projects", func(w http.ResponseWriter, r *http.Request) {
		_ = projects.Create(w, r)
	})

	return router, projectHandlerDeps{logger, projectService, columnService, taskService}
}

package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type projectService interface {
	List(ctx context.Context, all bool) ([]model.Project, error)
	Retrieve(ctx context.Context, projectID string) (model.Project, error)
	Create(ctx context.Context, project model.NewProject, now time.Time) (model.Project, error)
	Update(ctx context.Context, projectID string, update model.UpdateProject, now time.Time) (model.Project, error)
	Delete(ctx context.Context, projectID string) error
}

// ProjectHandler handles the project requests.
type ProjectHandler struct {
	logger         *zap.Logger
	projectService projectService
	columnService  columnService
	taskService    taskService
}

// NewProjectHandler returns a new project handler.
func NewProjectHandler(
	logger *zap.Logger,
	projectService projectService,
	columnService columnService,
	taskService taskService,
) *ProjectHandler {
	return &ProjectHandler{
		logger:         logger,
		projectService: projectService,
		columnService:  columnService,
		taskService:    taskService,
	}
}

// List handles project list requests.
func (ph *ProjectHandler) List(w http.ResponseWriter, r *http.Request) error {
	var all bool

	path := r.Header.Get("BasePath")
	if strings.ToLower(path) == "projects" {
		all = true
	}

	list, err := ph.projectService.List(r.Context(), all)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve handles project retrieval requests.
func (ph *ProjectHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	project, err := ph.projectService.Retrieve(r.Context(), pid)
	if err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error retrieving project %q: %w", pid, err)
		}
	}

	return web.Respond(r.Context(), w, project, http.StatusOK)
}

// Create handles project create requests.
func (ph *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var (
		np  model.NewProject
		err error
	)

	if err = web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	project, err := ph.projectService.Create(r.Context(), np, time.Now())
	if err != nil {
		return err
	}

	titles := [4]string{"To Do", "In Progress", "Review", "Done"}

	ph.logger.Info(project.ID)
	for i, title := range titles {
		nt := model.NewColumn{
			ProjectID:  project.ID,
			Title:      title,
			ColumnName: fmt.Sprintf(`column-%d`, i+1),
		}
		_, err = ph.columnService.Create(r.Context(), nt, time.Now())
		if err != nil {
			return err
		}
	}

	return web.Respond(r.Context(), w, project, http.StatusCreated)
}

// Update handles project update requests.
func (ph *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) error {
	var update model.UpdateProject

	pid := chi.URLParam(r, "pid")

	if err := web.Decode(r, &update); err != nil {
		return fmt.Errorf("error decoding project update: %w", err)
	}

	up, err := ph.projectService.Update(r.Context(), pid, update, time.Now())
	if err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error updating project %q: %w", pid, err)
		}
	}

	return web.Respond(r.Context(), w, up, http.StatusOK)
}

// Delete handles project delete requests.
func (ph *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	var err error
	pid := chi.URLParam(r, "pid")

	if _, err = ph.projectService.Retrieve(r.Context(), pid); err != nil {
		return err
	}
	if err = ph.taskService.DeleteAll(r.Context(), pid); err != nil {
		return err
	}
	if err = ph.columnService.DeleteAll(r.Context(), pid); err != nil {
		return err
	}
	if err = ph.projectService.Delete(r.Context(), pid); err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error deleting project %q: %w", pid, err)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

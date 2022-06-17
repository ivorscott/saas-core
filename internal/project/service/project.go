package service

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/project/model"

	"go.uber.org/zap"
)

type projectRepository interface {
	RetrieveTeamID(ctx context.Context, pid string) (string, error)
	Retrieve(ctx context.Context, pid string) (model.Project, error)
	RetrieveShared(ctx context.Context, pid string) (model.Project, error)
	List(ctx context.Context) ([]model.Project, error)
	Create(ctx context.Context, np model.NewProject, now time.Time) (model.Project, error)
	Update(ctx context.Context, pid string, update model.UpdateProject, now time.Time) (model.Project, error)
	Delete(ctx context.Context, pid string) error
}

// ProjectService is responsible for managing project business logic.
type ProjectService struct {
	logger *zap.Logger
	repo   projectRepository
}

// NewProjectService returns a new ProjectService.
func NewProjectService(logger *zap.Logger, repo projectRepository) *ProjectService {
	return &ProjectService{
		logger: logger,
		repo:   repo,
	}
}

// List lists projects.
func (ps *ProjectService) List(ctx context.Context) ([]model.Project, error) {
	return ps.repo.List(ctx)
}

// Retrieve retrieves an owned project.
func (ps *ProjectService) Retrieve(ctx context.Context, projectID string) (model.Project, error) {
	return ps.repo.Retrieve(ctx, projectID)
}

// RetrieveShared retrieves shared a project.
func (ps *ProjectService) RetrieveShared(ctx context.Context, projectID string) (model.Project, error) {
	return ps.repo.RetrieveShared(ctx, projectID)
}

// Create creates a project.
func (ps *ProjectService) Create(ctx context.Context, project model.NewProject, now time.Time) (model.Project, error) {
	return ps.repo.Create(ctx, project, now)
}

// Update updates a project.
func (ps *ProjectService) Update(ctx context.Context, projectID string, update model.UpdateProject, now time.Time) (model.Project, error) {
	return ps.repo.Update(ctx, projectID, update, now)
}

// Delete deletes a project.
func (ps *ProjectService) Delete(ctx context.Context, projectID string) error {
	return ps.repo.Delete(ctx, projectID)
}

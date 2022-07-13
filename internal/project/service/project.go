package service

import (
	"context"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/devpies/saas-core/internal/project/model"

	"go.uber.org/zap"
)

type projectRepository interface {
	RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error
	Retrieve(ctx context.Context, pid string) (model.Project, error)
	List(ctx context.Context) ([]model.Project, error)
	Create(ctx context.Context, np model.NewProject, now time.Time) (model.Project, error)
	Update(ctx context.Context, pid string, update model.UpdateProject, now time.Time) (model.Project, error)
	UpdateTx(ctx context.Context, tx *sqlx.Tx, pid string, update model.UpdateProject, now time.Time) (model.Project, error)
	Delete(ctx context.Context, pid string) error
}

// ProjectService is responsible for managing project business logic.
type ProjectService struct {
	logger      *zap.Logger
	projectRepo projectRepository
}

// NewProjectService returns a new ProjectService.
func NewProjectService(logger *zap.Logger, projectRepo projectRepository) *ProjectService {
	return &ProjectService{
		logger:      logger,
		projectRepo: projectRepo,
	}
}

// List retrieves projects across tenant accounts for the authenticated user.
func (ps *ProjectService) List(ctx context.Context, all bool) ([]model.Project, error) {
	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}
	if all {
		return forEachT(ctx, values.TenantMap, ps.projectRepo.List)
	}
	return ps.projectRepo.List(ctx)
}

// Retrieve retrieves an owned project.
func (ps *ProjectService) Retrieve(ctx context.Context, projectID string) (model.Project, error) {
	return ps.projectRepo.Retrieve(ctx, projectID)
}

// Create creates a project.
func (ps *ProjectService) Create(ctx context.Context, project model.NewProject, now time.Time) (model.Project, error) {
	return ps.projectRepo.Create(ctx, project, now)
}

// Update updates a project.
func (ps *ProjectService) Update(ctx context.Context, projectID string, update model.UpdateProject, now time.Time) (model.Project, error) {
	return ps.projectRepo.Update(ctx, projectID, update, now)
}

// Delete deletes a project.
func (ps *ProjectService) Delete(ctx context.Context, projectID string) error {
	return ps.projectRepo.Delete(ctx, projectID)
}

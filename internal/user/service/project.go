package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
)

type projectRepository interface {
	Create(ctx context.Context, p model.ProjectCopy) error
	Retrieve(ctx context.Context, pid string) (model.ProjectCopy, error)
	Update(ctx context.Context, pid string, update model.UpdateProjectCopy) error
	Delete(ctx context.Context, pid string) error
}

// ProjectService manages the redundant copy of project data.
type ProjectService struct {
	logger      *zap.Logger
	projectRepo projectRepository
}

// NewProjectService returns a new project service.
func NewProjectService(
	logger *zap.Logger,
	projectRepo projectRepository,
) *ProjectService {
	return &ProjectService{
		logger:      logger,
		projectRepo: projectRepo,
	}
}

func (ps *ProjectService) Create(ctx context.Context, p model.ProjectCopy) error {
	return ps.projectRepo.Create(ctx, p)
}

func (ps *ProjectService) Retrieve(ctx context.Context, pid string) (model.ProjectCopy, error) {
	return ps.projectRepo.Retrieve(ctx, pid)
}

func (ps *ProjectService) Update(ctx context.Context, pid string, update model.UpdateProjectCopy) error {
	return ps.projectRepo.Update(ctx, pid, update)
}

func (ps *ProjectService) Delete(ctx context.Context, pid string) error {
	return ps.projectRepo.Delete(ctx, pid)
}

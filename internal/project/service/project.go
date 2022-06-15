package service

import (
	"context"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"time"
)

type projectRepository interface{}

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

func (p ProjectService) List(ctx context.Context, userID string) ([]model.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectService) Retrieve(ctx context.Context, projectID string, userID string) (model.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectService) RetrieveShared(ctx context.Context, projectID string, userID string) (model.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectService) Create(ctx context.Context, project model.NewProject, userID string, now time.Time) (model.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectService) Update(ctx context.Context, projectID string, userID string, update model.UpdateProject, now time.Time) (model.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectService) Delete(ctx context.Context, projectID string, userID string) error {
	//TODO implement me
	panic("implement me")
}

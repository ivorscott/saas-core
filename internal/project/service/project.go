package service

import (
	"go.uber.org/zap"
)

type projectRepository interface{}

// ProjectService is responsible for managing project business logic.
type ProjectService struct {
	logger *zap.Logger
	js     publisher
	repo   projectRepository
}

// NewProjectService returns a new ProjectService.
func NewProjectService(logger *zap.Logger, js publisher, repo projectRepository) *ProjectService {
	return &ProjectService{
		logger: logger,
		js:     js,
		repo:   repo,
	}
}

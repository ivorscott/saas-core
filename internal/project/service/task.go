package service

import (
	"go.uber.org/zap"
)

type taskRepository interface{}

// TaskService is responsible for managing task business logic.
type TaskService struct {
	logger *zap.Logger
	js     publisher
	repo   taskRepository
}

// NewTaskService returns a TaskService.
func NewTaskService(logger *zap.Logger, js publisher, repo taskRepository) *TaskService {
	return &TaskService{
		logger: logger,
		js:     js,
		repo:   repo,
	}
}

package service

import (
	"context"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"time"
)

type taskRepository interface{}

// TaskService is responsible for managing task business logic.
type TaskService struct {
	logger *zap.Logger
	repo   taskRepository
}

// NewTaskService returns a TaskService.
func NewTaskService(logger *zap.Logger, repo taskRepository) *TaskService {
	return &TaskService{
		logger: logger,
		repo:   repo,
	}
}

func (t TaskService) Create(ctx context.Context, task model.NewTask, projectID string, userID string, now time.Time) (model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskService) List(ctx context.Context, projectID string) ([]model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskService) Retrieve(ctx context.Context, taskID string) (model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskService) Update(ctx context.Context, taskID string, update model.UpdateTask, now time.Time) (model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskService) Delete(ctx context.Context, taskID string) error {
	//TODO implement me
	panic("implement me")
}

func (t TaskService) DeleteAll(ctx context.Context, projectID string) error {
	//TODO implement me
	panic("implement me")
}

package service

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/project/model"

	"go.uber.org/zap"
)

type taskRepository interface {
	Create(ctx context.Context, nt model.NewTask, pid string, now time.Time) (model.Task, error)
	Retrieve(ctx context.Context, tid string) (model.Task, error)
	List(ctx context.Context, pid string) ([]model.Task, error)
	Update(ctx context.Context, tid string, update model.UpdateTask, now time.Time) (model.Task, error)
	Delete(ctx context.Context, tid string) error
}

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

// Create creates a task.
func (ts *TaskService) Create(ctx context.Context, task model.NewTask, projectID string, now time.Time) (model.Task, error) {
	return ts.repo.Create(ctx, task, projectID, now)
}

// List lists a task.
func (ts *TaskService) List(ctx context.Context, projectID string) ([]model.Task, error) {
	return ts.repo.List(ctx, projectID)
}

// Retrieve retrieves a task.
func (ts *TaskService) Retrieve(ctx context.Context, taskID string) (model.Task, error) {
	return ts.repo.Retrieve(ctx, taskID)
}

// Update updates a task.
func (ts *TaskService) Update(ctx context.Context, taskID string, update model.UpdateTask, now time.Time) (model.Task, error) {
	return ts.repo.Update(ctx, taskID, update, now)
}

// Delete deletes a task.
func (ts *TaskService) Delete(ctx context.Context, taskID string) error {
	return ts.repo.Delete(ctx, taskID)
}

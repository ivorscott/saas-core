package handler

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/project/model"
)

type columnService interface {
	Create(ctx context.Context, nc model.NewColumn, now time.Time) (model.Column, error)
	CreateColumns(ctx context.Context, pid string, now time.Time) error
	List(ctx context.Context, projectID string) ([]model.Column, error)
	Retrieve(ctx context.Context, columnID string) (model.Column, error)
	Update(ctx context.Context, columnID string, update model.UpdateColumn, now time.Time) (model.Column, error)
	Delete(ctx context.Context, columnID string) error
}

type taskService interface {
	Create(ctx context.Context, task model.NewTask, projectID string, now time.Time) (model.Task, error)
	List(ctx context.Context, projectID string) ([]model.Task, error)
	Retrieve(ctx context.Context, taskID string) (model.Task, error)
	Update(ctx context.Context, taskID string, update model.UpdateTask, now time.Time) (model.Task, error)
	Delete(ctx context.Context, taskID string) error
}

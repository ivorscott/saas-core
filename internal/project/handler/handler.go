package handler

import (
	"context"
	"github.com/devpies/saas-core/internal/project/model"
	"time"
)

type columnService interface {
	Create(ctx context.Context, column model.NewColumn, now time.Time) (model.Column, error)
	List(ctx context.Context, projectID string) ([]model.Column, error)
	Retrieve(ctx context.Context, columnID string) (model.Column, error)
	Update(ctx context.Context, columnID string, update model.UpdateColumn, now time.Time) error
	Delete(ctx context.Context, columnID string) error
	DeleteAll(ctx context.Context, projectID string) error
}

type taskService interface {
	Create(ctx context.Context, task model.NewTask, projectID string, userID string, now time.Time) (model.Task, error)
	List(ctx context.Context, projectID string) ([]model.Task, error)
	Retrieve(ctx context.Context, taskID string) (model.Task, error)
	Update(ctx context.Context, taskID string, update model.UpdateTask, now time.Time) (model.Task, error)
	Delete(ctx context.Context, taskID string) error
	DeleteAll(ctx context.Context, projectID string) error
}

type publisher interface {
	Publish(subject string, message []byte)
}

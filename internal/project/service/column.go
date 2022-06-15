package service

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/project/model"

	"go.uber.org/zap"
)

type columnRepository interface {
	Create(ctx context.Context, nc model.NewColumn, now time.Time) (model.Column, error)
	Retrieve(ctx context.Context, cid string) (model.Column, error)
	List(ctx context.Context, pid string) ([]model.Column, error)
	Update(ctx context.Context, cid string, uc model.UpdateColumn, now time.Time) error
	Delete(ctx context.Context, cid string) error
	DeleteAll(ctx context.Context, pid string) error
}

// ColumnService is responsible for managing column business logic.
type ColumnService struct {
	logger *zap.Logger
	repo   columnRepository
}

// NewColumnService returns a new ColumnService.
func NewColumnService(logger *zap.Logger, repo columnRepository) *ColumnService {
	return &ColumnService{
		logger: logger,
		repo:   repo,
	}
}

// Create creates a new project column.
func (cs *ColumnService) Create(ctx context.Context, column model.NewColumn, now time.Time) (model.Column, error) {
	return cs.repo.Create(ctx, column, now)
}

// List lists all project columns.
func (cs *ColumnService) List(ctx context.Context, projectID string) ([]model.Column, error) {
	return cs.repo.List(ctx, projectID)
}

// Retrieve retrieves a project column.
func (cs *ColumnService) Retrieve(ctx context.Context, columnID string) (model.Column, error) {
	return cs.repo.Retrieve(ctx, columnID)
}

// Update updates a project column.
func (cs *ColumnService) Update(ctx context.Context, columnID string, update model.UpdateColumn, now time.Time) error {
	return cs.repo.Update(ctx, columnID, update, now)
}

// Delete deletes a project column.
func (cs *ColumnService) Delete(ctx context.Context, columnID string) error {
	return cs.repo.Delete(ctx, columnID)
}

// DeleteAll deletes all project columns.
func (cs *ColumnService) DeleteAll(ctx context.Context, projectID string) error {
	return cs.repo.DeleteAll(ctx, projectID)
}

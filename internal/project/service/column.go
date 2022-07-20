package service

import (
	"context"
	"fmt"
	"time"

	"github.com/devpies/saas-core/internal/project/model"

	"go.uber.org/zap"
)

type columnRepository interface {
	Create(ctx context.Context, nc model.NewColumn, now time.Time) (model.Column, error)
	Retrieve(ctx context.Context, cid string) (model.Column, error)
	List(ctx context.Context, pid string) ([]model.Column, error)
	Update(ctx context.Context, cid string, uc model.UpdateColumn, now time.Time) (model.Column, error)
	Delete(ctx context.Context, cid string) error
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

// CreateColumns creates new project columns.
func (cs *ColumnService) CreateColumns(ctx context.Context, pid string, now time.Time) error {
	titles := [4]string{"To Do", "In Progress", "Review", "Done"}

	for i, title := range titles {
		nc := model.NewColumn{
			ProjectID:  pid,
			Title:      title,
			ColumnName: fmt.Sprintf(`column-%d`, i+1),
		}
		_, err := cs.repo.Create(ctx, nc, now)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create creates a new project column.
func (cs *ColumnService) Create(ctx context.Context, nc model.NewColumn, now time.Time) (model.Column, error) {
	column, err := cs.repo.Create(ctx, nc, now)
	if err != nil {
		return model.Column{}, err
	}
	return column, nil
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
func (cs *ColumnService) Update(ctx context.Context, columnID string, update model.UpdateColumn, now time.Time) (model.Column, error) {
	return cs.repo.Update(ctx, columnID, update, now)
}

// Delete deletes a project column.
func (cs *ColumnService) Delete(ctx context.Context, columnID string) error {
	return cs.repo.Delete(ctx, columnID)
}

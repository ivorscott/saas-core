package service

import (
	"context"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"time"
)

type columnRepository interface{}

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

func (c ColumnService) Create(ctx context.Context, column model.NewColumn, now time.Time) (model.Column, error) {
	//TODO implement me
	panic("implement me")
}

func (c ColumnService) List(ctx context.Context, projectID string) ([]model.Column, error) {
	//TODO implement me
	panic("implement me")
}

func (c ColumnService) Retrieve(ctx context.Context, columnID string) (model.Column, error) {
	//TODO implement me
	panic("implement me")
}

func (c ColumnService) Update(ctx context.Context, columnID string, update model.UpdateColumn, now time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c ColumnService) Delete(ctx context.Context, columnID string) error {
	//TODO implement me
	panic("implement me")
}

func (c ColumnService) DeleteAll(ctx context.Context, projectID string) error {
	//TODO implement me
	panic("implement me")
}

package service

import (
	"go.uber.org/zap"
)

type columnRepository interface{}

// ColumnService is responsible for managing column business logic.
type ColumnService struct {
	logger *zap.Logger
	js     publisher
	repo   columnRepository
}

// NewColumnService returns a new ColumnService.
func NewColumnService(logger *zap.Logger, js publisher, repo columnRepository) *ColumnService {
	return &ColumnService{
		logger: logger,
		js:     js,
		repo:   repo,
	}
}

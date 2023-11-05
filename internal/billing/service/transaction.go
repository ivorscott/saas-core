package service

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/billing/model"

	"go.uber.org/zap"
)

type transactionRepository interface {
	SaveTransaction(ctx context.Context, nt model.NewTransaction, now time.Time) (model.Transaction, error)
}

// TransactionService is responsible for managing transaction related business logic.
type TransactionService struct {
	logger *zap.Logger
	repo   transactionRepository
}

// NewTransactionService returns a new TransactionService.
func NewTransactionService(logger *zap.Logger, repo transactionRepository) *TransactionService {
	return &TransactionService{
		logger: logger,
		repo:   repo,
	}
}

// Save persists the new transaction details.
func (ts *TransactionService) Save(ctx context.Context, nt model.NewTransaction, now time.Time) (model.Transaction, error) {
	var (
		t   model.Transaction
		err error
	)

	t, err = ts.repo.SaveTransaction(ctx, nt, now)
	if err != nil {
		switch err {
		default:
			return t, err
		}
	}
	return t, nil
}

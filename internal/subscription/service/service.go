package service

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/subscription/model"
)

type customerRepository interface {
	GetCustomer(ctx context.Context, tenantID string) (model.Customer, error)
	SaveCustomer(ctx context.Context, nc model.NewCustomer, now time.Time) (model.Customer, error)
}

type transactionRepository interface {
	SaveTransaction(ctx context.Context, nt model.NewTransaction, now time.Time) (model.Transaction, error)
	GetAllTransactions(ctx context.Context, tenantID string) ([]model.Transaction, error)
}

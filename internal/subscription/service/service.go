package service

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/devpies/saas-core/internal/subscription/model"
)

type customerRepository interface {
	GetCustomer(ctx context.Context, tenantID string) (model.Customer, error)
	SaveCustomerTx(ctx context.Context, tx *sqlx.Tx, nc model.NewCustomer, now time.Time) (model.Customer, error)
	RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error
}

type transactionRepository interface {
	SaveTransactionTx(ctx context.Context, tx *sqlx.Tx, nt model.NewTransaction, now time.Time) (model.Transaction, error)
	GetAllTransactions(ctx context.Context, tenantID string) ([]model.Transaction, error)
}

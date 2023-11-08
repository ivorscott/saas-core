package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/devpies/saas-core/internal/subscription/db"
	"github.com/devpies/saas-core/internal/subscription/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TransactionRepository manages data access to customer transactions.
type TransactionRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewTransactionRepository returns a new TransactionRepository.
func NewTransactionRepository(logger *zap.Logger, pg *db.PostgresDatabase) *TransactionRepository {
	return &TransactionRepository{
		logger: logger,
		pg:     pg,
	}
}

// SaveTransaction saves a new customer transaction.
func (tr *TransactionRepository) SaveTransaction(ctx context.Context, nt model.NewTransaction, now time.Time) (model.Transaction, error) {
	var (
		t   model.Transaction
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return t, web.CtxErr()
	}
	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, err
	}
	defer Close()

	t = model.Transaction{
		ID:                   uuid.New().String(),
		Amount:               nt.Amount,
		Currency:             nt.Currency,
		LastFour:             nt.LastFour,
		BankReturnCode:       nt.BankReturnCode,
		StatusID:             nt.StatusID,
		ExpirationMonth:      nt.ExpirationMonth,
		ExpirationYear:       nt.ExpirationYear,
		StripeSubscriptionID: nt.StripeSubscriptionID,
		PaymentIntent:        nt.PaymentIntent,
		PaymentMethod:        nt.PaymentMethod,
		TenantID:             values.TenantID,
		UpdatedAt:            now.UTC(),
		CreatedAt:            now.UTC(),
	}

	stmt := `
			insert into transactions (
				transaction_id, amount, currency, last_four, bank_return_code,
				transaction_status_id, expiration_month, expiration_year, stripe_subscription_id,
				payment_intent, payment_method, tenant_id, updated_at, created_at
			) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`
	if _, err = conn.ExecContext(
		ctx,
		stmt,
		t.ID,
		t.Amount,
		t.Currency,
		t.LastFour,
		t.BankReturnCode,
		t.StatusID,
		t.ExpirationMonth,
		t.ExpirationYear,
		t.StripeSubscriptionID,
		t.PaymentIntent,
		t.PaymentMethod,
		t.TenantID,
		t.UpdatedAt,
		t.CreatedAt,
	); err != nil {
		return t, fmt.Errorf("error inserting transaction %v :%w", t, err)
	}

	return t, nil
}

// GetAllTransactions retrieves all transactions for the provided tenant.
func (tr *TransactionRepository) GetAllTransactions(ctx context.Context, tenantID string) ([]model.Transaction, error) {
	var (
		t   model.Transaction
		ts  = make([]model.Transaction, 0)
		err error
	)

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return ts, err
	}
	defer Close()

	stmt := `
			select transaction_id, amount, currency, last_four, bank_return_code, transaction_status_id,
				expiration_month, expiration_year, stripe_subscription_id, payment_intent, payment_method,
				tenant_id, updated_at, created_at
			from transactions
			where tenant_id = $1
	`

	rows, err := conn.QueryxContext(ctx, stmt, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error selecting transactions :%w", err)
	}
	for rows.Next() {
		err = rows.Scan(
			&t.ID,
			&t.Amount,
			&t.Currency,
			&t.LastFour,
			&t.BankReturnCode,
			&t.StatusID,
			&t.ExpirationMonth,
			&t.ExpirationYear,
			&t.StripeSubscriptionID,
			&t.PaymentIntent,
			&t.PaymentMethod,
			&t.TenantID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row into struct :%w", err)
		}

		t.CreatedAt = t.CreatedAt.UTC()
		t.UpdatedAt = t.UpdatedAt.UTC()

		ts = append(ts, t)
	}

	return ts, nil
}

// Package repository manages the data access layer for handling queries.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/devpies/saas-core/internal/subscription/db"
	"github.com/devpies/saas-core/internal/subscription/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SubscriptionRepository manages data access to subscriptions.
type SubscriptionRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

var (
	// ErrNotFound represents a resource not found.
	ErrNotFound = errors.New("not found")
)

// NewSubscriptionRepository returns a new SubscriptionRepository.
func NewSubscriptionRepository(logger *zap.Logger, pg *db.PostgresDatabase) *SubscriptionRepository {
	return &SubscriptionRepository{
		logger: logger,
		pg:     pg,
	}
}

// SaveSubscription saves a subscription.
func (sr *SubscriptionRepository) SaveSubscription(ctx context.Context, ns model.NewSubscription, now time.Time) (model.Subscription, error) {
	var (
		s   model.Subscription
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return s, web.CtxErr()
	}
	conn, Close, err := sr.pg.GetConnection(ctx)
	if err != nil {
		return s, err
	}
	defer Close()

	stmt := `
			insert into subscriptions (
   				subscription_id, plan, transaction_id, subscription_status_id,
				amount, customer_id, tenant_id
			) values ($1, $2, $3, $4, $5, $6, $7)
		`
	s = model.Subscription{
		ID:            uuid.New().String(),
		Plan:          ns.Plan,
		TransactionID: ns.TransactionID,
		StatusID:      ns.StatusID,
		Amount:        ns.Amount,
		TenantID:      values.TenantID,
		CustomerID:    ns.CustomerID,
		UpdatedAt:     now.UTC(),
		CreatedAt:     now.UTC(),
	}

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		s.ID,
		s.Plan,
		s.TransactionID,
		s.StatusID,
		s.Amount,
		s.CustomerID,
		s.TenantID,
	); err != nil {
		return s, fmt.Errorf("error inserting subscription %v :%w", s, err)
	}

	return s, nil
}

// GetTenantSubscription retrieves the tenant subscription by tenant id.
// Because the admin service will also make calls to the billing service,
// we cannot rely on there being a tenant id in the request context.
func (sr *SubscriptionRepository) GetTenantSubscription(ctx context.Context, tenantID string) (model.Subscription, error) {
	var (
		s   model.Subscription
		err error
	)

	conn, Close, err := sr.pg.GetConnection(ctx)
	if err != nil {
		return s, err
	}
	defer Close()

	stmt := `
			select subscription_id, tenant_id, customer_id, transaction_id,
					subscription_status_id, amount, plan, updated_at, created_at
			from subscriptions
			where tenant_id = $1
		`
	row := conn.QueryRowxContext(ctx, stmt, tenantID)

	if err = row.Scan(
		&s.ID,
		&s.TenantID,
		&s.CustomerID,
		&s.TransactionID,
		&s.StatusID,
		&s.Amount,
		&s.Plan,
		&s.UpdatedAt,
		&s.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return s, ErrNotFound
		}
		return s, err
	}

	return s, nil
}

package repository

import (
	"context"
	"fmt"
	"github.com/devpies/saas-core/internal/billing/db"
	"github.com/devpies/saas-core/internal/billing/model"
	"github.com/devpies/saas-core/pkg/web"
	"go.uber.org/zap"
	"time"
)

type SubscriptionRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

func NewSubscriptionRepository(logger *zap.Logger, pg *db.PostgresDatabase) *SubscriptionRepository {
	return &SubscriptionRepository{
		logger: logger,
		pg:     pg,
	}
}

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
		ID:            ns.ID,
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

func (sr *SubscriptionRepository) GetAllSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	var subs []model.Subscription
	return subs, nil
}
func (sr *SubscriptionRepository) GetOneSubscription(ctx context.Context, id string) (model.Subscription, error) {
	var s model.Subscription
	return s, nil
}

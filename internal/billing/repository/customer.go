package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/devpies/saas-core/internal/billing/db"
	"github.com/devpies/saas-core/internal/billing/model"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

// CustomerRepository manages data access to customers.
type CustomerRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewCustomerRepository returns a CustomerRepository.
func NewCustomerRepository(logger *zap.Logger, pg *db.PostgresDatabase) *CustomerRepository {
	return &CustomerRepository{
		logger: logger,
		pg:     pg,
	}
}

// SaveCustomer saves a customer.
func (cr *CustomerRepository) SaveCustomer(ctx context.Context, nc model.NewCustomer, now time.Time) (model.Customer, error) {
	var c model.Customer

	values, ok := web.FromContext(ctx)
	if !ok {
		return c, web.CtxErr()
	}
	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return c, err
	}
	defer Close()

	c = model.Customer{
		ID:        nc.ID,
		TenantID:  values.TenantID,
		FirstName: nc.FirstName,
		LastName:  nc.LastName,
		Email:     nc.Email,
		UpdatedAt: now.UTC(),
		CreatedAt: now.UTC(),
	}

	stmt := `
			insert into customers (customer_id, tenant_id, first_name, last_name, email)
			values ($1,$2,$3,$4,$5)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		c.ID,
		c.TenantID,
		c.FirstName,
		c.LastName,
		c.Email,
	); err != nil {
		return c, fmt.Errorf("error inserting customer %v :%w", c, err)
	}

	return c, nil
}

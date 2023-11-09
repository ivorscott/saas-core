package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/devpies/saas-core/internal/subscription/db"
	"github.com/devpies/saas-core/internal/subscription/model"
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

	if err := nc.Validate(); err != nil {
		return c, err
	}

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
			insert into customers (customer_id, tenant_id, first_name, last_name, email, updated_at, created_at)
			values ($1,$2,$3,$4,$5,$6,$7)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		c.ID,
		c.TenantID,
		c.FirstName,
		c.LastName,
		c.Email,
		c.UpdatedAt,
		c.CreatedAt,
	); err != nil {
		return c, fmt.Errorf("error inserting customer %v :%w", c, err)
	}

	return c, nil
}

// GetCustomer retrieves the customer.
func (cr *CustomerRepository) GetCustomer(ctx context.Context, tenantID string) (model.Customer, error) {
	var (
		c   model.Customer
		err error
	)
	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return c, err
	}
	defer Close()
	stmt := `
			select customer_id, tenant_id, first_name, last_name, email,
				updated_at, created_at
			from customers
			where tenant_id = $1
	`

	row := conn.QueryRowContext(ctx, stmt, tenantID)
	if err = row.Scan(
		&c.ID,
		&c.TenantID,
		&c.FirstName,
		&c.LastName,
		&c.Email,
		&c.UpdatedAt,
		&c.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return c, ErrNotFound
		}
		return c, err
	}
	return c, nil
}

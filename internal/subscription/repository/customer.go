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

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// CustomerRepository manages data access to customers.
type CustomerRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
	runTx  func(ctx context.Context, fn func(*sqlx.Tx) error) error
}

var (
	// ErrCustomerNotFound represents a customer not found.
	ErrCustomerNotFound = errors.New("customer not found")
)

// NewCustomerRepository returns a CustomerRepository.
func NewCustomerRepository(logger *zap.Logger, pg *db.PostgresDatabase) *CustomerRepository {
	return &CustomerRepository{
		logger: logger,
		pg:     pg,
		runTx:  pg.RunInTransaction,
	}
}

// RunTx runs a function within a transaction context.
func (cr *CustomerRepository) RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	return cr.runTx(ctx, fn)
}

// SaveCustomerTx saves a customer.
func (cr *CustomerRepository) SaveCustomerTx(ctx context.Context, tx *sqlx.Tx, nc model.NewCustomer, now time.Time) (model.Customer, error) {
	var (
		c   model.Customer
		err error
	)

	if err := nc.Validate(); err != nil {
		return c, err
	}

	values, ok := web.FromContext(ctx)
	if !ok {
		return c, web.CtxErr()
	}

	c = model.Customer{
		ID:              nc.ID,
		TenantID:        values.TenantID,
		FirstName:       nc.FirstName,
		LastName:        nc.LastName,
		Email:           nc.Email,
		PaymentMethodID: nc.PaymentMethodID,
		UpdatedAt:       now.UTC(),
		CreatedAt:       now.UTC(),
	}

	stmt := `
			insert into customers (customer_id, tenant_id, first_name, last_name, email, payment_method)
			values ($1,$2,$3,$4,$5,$6)
	`

	if _, err = tx.ExecContext(
		ctx,
		stmt,
		c.ID,
		c.TenantID,
		c.FirstName,
		c.LastName,
		c.Email,
		c.PaymentMethodID,
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
			select customer_id, tenant_id, first_name, last_name, email, payment_method,
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
		&c.PaymentMethodID,
		&c.UpdatedAt,
		&c.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return c, ErrCustomerNotFound
		}
		return c, err
	}
	return c, nil
}

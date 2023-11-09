// Package repository manages the data access layer for handling queries.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"time"

	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// UserRepository manages user data access.
type UserRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
	runTx  func(ctx context.Context, fn func(*sqlx.Tx) error) error
}

// NewUserRepository returns a new user repository.
func NewUserRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *UserRepository {
	return &UserRepository{
		logger: logger,
		pg:     pg,
		runTx:  pg.RunInTransaction,
	}
}

// RunTx runs a function within a transaction context.
func (ur *UserRepository) RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	return ur.runTx(ctx, fn)
}

// AddUserTx adds a user to a tenant in the database.
func (ur *UserRepository) AddUserTx(ctx context.Context, tx *sqlx.Tx, userID string, now time.Time) error {
	var (
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}

	stmt := `
		insert into users (user_id, tenant_id, created_at)
		values ($1, $2, $3)
	`

	if _, err = tx.ExecContext(
		ctx,
		stmt,
		userID,
		values.TenantID,
		now.UTC(),
	); err != nil {
		return err
	}

	return nil
}

// CreateUserProfile creates a tenant agnostic user profile. Row level security is not enabled on "user_profiles".
func (ur *UserRepository) CreateUserProfile(ctx context.Context, nu model.NewUser, userID string, now time.Time) (model.User, error) {
	var (
		u   model.User
		err error
	)

	conn, Close, err := ur.pg.GetConnection(ctx)
	if err != nil {
		return u, fail.ErrConnectionFailed
	}
	defer Close()

	u = model.User{
		ID:        userID,
		Email:     nu.Email,
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		UpdatedAt: now.UTC(),
		CreatedAt: now.UTC(),
	}

	stmt := `
		insert into user_profiles (user_id, email, email_verified, first_name, last_name, picture, locale, updated_at, created_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.Email,
		u.EmailVerified,
		u.FirstName,
		u.LastName,
		u.Picture,
		u.Locale,
		u.UpdatedAt,
		u.CreatedAt,
	); err != nil {
		return u, err
	}

	return u, nil
}

// CreateAdminUser creates an admin user for the tenant.
func (ur *UserRepository) CreateAdminUser(ctx context.Context, na model.NewAdminUser) error {
	var err error

	conn, Close, err := ur.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		insert into user_profiles (user_id, email, email_verified, first_name, last_name, picture, locale, updated_at, created_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		na.UserID,
		na.Email,
		na.EmailVerified,
		na.FirstName,
		na.LastName,
		nil,
		nil,
		na.CreatedAt.UTC(),
		na.CreatedAt.UTC(),
	); err != nil {
		return err
	}

	stmt = `
		insert into users (user_id, tenant_id, created_at)
		values ($1, $2, $3)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		na.UserID,
		na.TenantID,
		na.CreatedAt.UTC(),
	); err != nil {
		return err
	}

	return nil
}

// List selects all users associated to the tenant account.
func (ur *UserRepository) List(ctx context.Context) ([]model.User, error) {
	var (
		u   model.User
		us  = make([]model.User, 0)
		err error
	)

	conn, Close, err := ur.pg.GetConnection(ctx)
	if err != nil {
		return us, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `select * from users inner join user_profiles using(user_id)`

	rows, err := conn.QueryxContext(ctx, stmt)
	if err != nil {
		return us, err
	}
	for rows.Next() {
		err = rows.StructScan(&u)
		if err != nil {
			return us, fmt.Errorf("error decoding struct: %w", err)
		}
		us = append(us, u)
	}

	return us, nil
}

// RetrieveIDByEmail retrieves a userID via a provided email address.
func (ur *UserRepository) RetrieveIDByEmail(ctx context.Context, email string) (string, error) {
	var (
		uid string
		err error
	)

	if _, err = mail.ParseAddress(email); err != nil {
		return "", fail.ErrInvalidEmail
	}

	conn, Close, err := ur.pg.GetConnection(ctx)
	if err != nil {
		return "", fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `select user_id from user_profiles where email = $1 limit 1`

	if err = conn.GetContext(ctx, &uid, stmt, email); err != nil {
		return "", err
	}

	return uid, nil
}

// RetrieveByEmail retrieves a user via a provided email address.
func (ur *UserRepository) RetrieveByEmail(ctx context.Context, email string) (model.User, error) {
	var (
		u   model.User
		err error
	)

	if _, err = mail.ParseAddress(email); err != nil {
		return u, fail.ErrInvalidEmail
	}

	conn, Close, err := ur.pg.GetConnection(ctx)
	if err != nil {
		return u, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			select 
			    u.user_id, tenant_id, email, first_name, last_name,
			    email_verified, locale, picture, u.created_at
			from users u
			inner join user_profiles using (user_id)
			where email = $1
			limit 1
	`

	if err = conn.GetContext(ctx, &u, stmt, email); err != nil {
		return u, err
	}

	return u, nil
}

// RetrieveMe retrieves the authenticated user.
func (ur *UserRepository) RetrieveMe(ctx context.Context) (model.User, error) {
	var (
		u   model.User
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return u, web.CtxErr()
	}

	if _, err = uuid.Parse(values.UserID); err != nil {
		return u, fail.ErrInvalidID
	}

	conn, Close, err := ur.pg.GetConnection(ctx)
	if err != nil {
		return u, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select 
			u.user_id, u.tenant_id, email, first_name, last_name,
			email_verified, locale, picture, u.created_at
		from users u
		inner join user_profiles using (user_id)
		where user_id = $1
	`

	if err = conn.GetContext(ctx, &u, stmt, values.UserID); err != nil {
		if err == sql.ErrNoRows {
			return u, fail.ErrNotFound
		}
		return u, err
	}

	return u, nil
}

// DetachUserTx detaches a user from a tenant account.
func (ur *UserRepository) DetachUserTx(ctx context.Context, tx *sqlx.Tx, uid string) error {
	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}

	stmt := `delete from users where user_id = $1 and tenant_id = $2`

	_, err := tx.ExecContext(ctx, stmt, uid, values.TenantID)
	if err != nil {
		ur.logger.Error("failed to remove user", zap.Error(err))
		return err
	}

	return nil
}

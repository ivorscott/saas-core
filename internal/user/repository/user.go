package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/mail"
	"time"

	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/google/uuid"
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

func (ur *UserRepository) RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	return ur.runTx(ctx, fn)
}

// Create inserts a new user into the database.
func (ur *UserRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, nu model.NewUser, now time.Time) (model.User, error) {
	var (
		u   model.User
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return u, web.CtxErr()
	}

	u = model.User{
		ID:            uuid.New().String(),
		TenantID:      values.TenantID,
		Email:         nu.Email,
		EmailVerified: nu.EmailVerified,
		FirstName:     nu.FirstName,
		LastName:      nu.LastName,
		Picture:       nu.Picture,
		Locale:        nu.Locale,
		UpdatedAt:     now.UTC(),
		CreatedAt:     now.UTC(),
	}

	stmt := `
		insert into users (user_id, tenant_id, email, email_verified, first_name, last_name, picture, locale, updated_at, created_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	if _, err = tx.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.TenantID,
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

// CreateAdmin inserts a new tenant admin user into the database.
func (ur *UserRepository) CreateAdminTx(ctx context.Context, tx *sqlx.Tx, na model.NewAdminUser) error {
	var (
		u   model.User
		err error
	)

	u = model.User{
		ID:            na.UserID,
		TenantID:      na.TenantID,
		Email:         na.Email,
		EmailVerified: na.EmailVerified,
		FirstName:     na.FirstName,
		LastName:      na.LastName,
		UpdatedAt:     na.CreatedAt.UTC(),
		CreatedAt:     na.CreatedAt.UTC(),
	}

	stmt := `
		insert into users (user_id, tenant_id, email, email_verified, first_name, last_name, picture, locale, updated_at, created_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	if _, err = tx.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.TenantID,
		u.Email,
		u.EmailVerified,
		u.FirstName,
		u.LastName,
		u.Picture,
		u.Locale,
		u.UpdatedAt,
		u.CreatedAt,
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

	stmt := `select * from users`

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
			    user_id, tenant_id, email, first_name, last_name,
			    email_verified, locale, picture, updated_at, created_at
			from users
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
			user_id, tenant_id, email, first_name, last_name,
			email_verified, locale, picture, updated_at, created_at
		from users
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

func (ur *UserRepository) RemoveUserTx(ctx context.Context, tx *sqlx.Tx, uid string) error {
	stmt := `delete from users where user_id = $1`

	_, err := tx.ExecContext(ctx, stmt, uid)
	if err != nil {
		ur.logger.Error("failed to remove user", zap.Error(err))
		return err
	}

	return nil
}

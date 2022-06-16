package repository

import (
	"context"
	"database/sql"
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
}

// NewUserRepository returns a new user repository.
func NewUserRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *UserRepository {
	return &UserRepository{
		logger: logger,
		pg:     pg,
	}
}

// Create inserts a new user into the database.
func (ur *UserRepository) Create(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error) {
	var (
		u   model.User
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return u, web.CtxErr()
	}

	conn, Close, err := ur.pg.GetConnection(ctx)
	if err != nil {
		return u, fail.ErrConnectionFailed
	}
	defer Close()

	u = model.User{
		ID:            uuid.New().String(),
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
		insert into users (user_id, tenant_id, email, email_verified, first_name, last_name, picture, local, update_at, created_at)
		values (?,?,?,?,?,?,?,?,?,?)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		u.ID,
		values.Metadata.TenantID,
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
			    email_verification, locale, picture, updated_at, created_at
			from users
			where email = ?
	`

	if err = conn.SelectContext(ctx, &u, stmt, email); err != nil {
		if err == sql.ErrNoRows {
			return u, fail.ErrNotFound
		}
		return u, err
	}

	return u, nil
}

// RetrieveMe retrieves the authenticated user.
func (ur *UserRepository) RetrieveMe(ctx context.Context, uid string) (model.User, error) {
	var (
		u   model.User
		err error
	)

	if _, err = uuid.Parse(uid); err != nil {
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
			    email_verification, locale, picture, updated_at, created_at
			from users
			where user_id = ?
	`

	if err = conn.SelectContext(ctx, &u, stmt, uid); err != nil {
		if err == sql.ErrNoRows {
			return u, fail.ErrNotFound
		}
		return u, err
	}

	return u, nil
}

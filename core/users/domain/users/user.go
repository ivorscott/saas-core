package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/google/uuid"
)

// Error codes returned by failures to handle users.
var (
	ErrNotFound     = errors.New("user not found")
	ErrInvalidID    = errors.New("id provided was not a valid UUID")
	ErrInvalidEmail = errors.New("address provided was not a valid email")
)

type UserQuerier interface {
	Create(ctx context.Context, repo database.Storer, nu NewUser, now time.Time) (User, error)
	RetrieveByEmail(repo database.Storer, email string) (User, error)
	RetrieveMe(ctx context.Context, repo database.Storer, uid string) (User, error)
	RetrieveMeByAuthID(ctx context.Context, repo database.Storer, aid string) (User, error)
}

type Queries struct{}

func (q *Queries) Create(ctx context.Context, repo database.Storer, nu NewUser, now time.Time) (User, error) {
	u := User{
		ID:            uuid.New().String(),
		Auth0ID:       nu.Auth0ID,
		Email:         nu.Email,
		EmailVerified: nu.EmailVerified,
		FirstName:     nu.FirstName,
		LastName:      nu.LastName,
		Picture:       nu.Picture,
		Locale:        nu.Locale,
		UpdatedAt:     now.UTC(),
		CreatedAt:     now.UTC(),
	}

	stmt := repo.Insert(
		"users",
	).SetMap(map[string]interface{}{
		"user_id":        u.ID,
		"auth0_id":       u.Auth0ID,
		"email":          u.Email,
		"email_verified": u.EmailVerified,
		"first_name":     u.FirstName,
		"last_name":      u.LastName,
		"picture":        u.Picture,
		"locale":         u.Locale,
		"updated_at":     u.UpdatedAt,
		"created_at":     u.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return u, err
	}

	return u, nil
}

func (q *Queries) RetrieveByEmail(repo database.Storer, email string) (User, error) {
	var u User

	if _, err := mail.ParseAddress(email); err != nil {
		return u, ErrInvalidEmail
	}

	stmt := repo.Select(
		"user_id",
		"auth0_id",
		"email",
		"first_name",
		"last_name",
		"email_verified",
		"locale",
		"picture",
		"updated_at",
		"created_at",
	).From(
		"users",
	).Where(sq.Eq{"email": "?"})

	query, args, err := stmt.ToSql()
	if err != nil {
		return u, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	if err := repo.Get(&u, query, email); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}

	return u, nil
}

func (q *Queries) RetrieveMe(ctx context.Context, repo database.Storer, uid string) (User, error) {
	var u User

	if _, err := uuid.Parse(uid); err != nil {
		return u, ErrInvalidID
	}
	stmt := repo.Select(
		"user_id",
		"auth0_id",
		"email",
		"first_name",
		"last_name",
		"email_verified",
		"locale",
		"picture",
		"updated_at",
		"created_at",
	).From(
		"users",
	).Where(sq.Eq{"user_id": "?"})

	query, args, err := stmt.ToSql()
	if err != nil {
		return u, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	if err := repo.GetContext(ctx, &u, query, uid); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}

	return u, nil
}

func (q *Queries) RetrieveMeByAuthID(ctx context.Context, repo database.Storer, aid string) (User, error) {
	var u User

	stmt := repo.Select(
		"user_id",
		"auth0_id",
		"email",
		"first_name",
		"last_name",
		"email_verified",
		"locale",
		"picture",
		"updated_at",
		"created_at",
	).From(
		"users",
	).Where(sq.Eq{"auth0_id": "?"})

	query, args, err := stmt.ToSql()
	if err != nil {
		return u, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	if err := repo.GetContext(ctx, &u, query, aid); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}

	return u, nil
}

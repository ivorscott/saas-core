package users

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrNotFound  = errors.New("user not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Create(ctx context.Context, repo *database.Repository, nu NewUser, aid string, now time.Time) (User, error) {
	u := User{
		ID:            uuid.New().String(),
		Auth0ID:       aid,
		Email:         nu.Email,
		EmailVerified: nu.EmailVerified,
		FirstName:     nu.FirstName,
		LastName:      nu.LastName,
		Picture:       nu.Picture,
		Locale:        nu.Locale,
		UpdatedAt: now.UTC(),
		CreatedAt: now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"users",
	).SetMap(map[string]interface{}{
		"user_id":        u.ID,
		"auth0_id":       u.Auth0ID, // unique
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
		return u, errors.Wrapf(err, "inserting user: %v", nu)
	}

	return u, nil
}

func RetrieveByEmail(repo *database.Repository, email string) (User, error) {
	var u User

	stmt := repo.SQ.Select(
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

	q, args, err := stmt.ToSql()
	if err != nil {
		return u, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.Get(&u, q, email); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}

	return u, nil
}

func RetrieveMe(ctx context.Context, repo *database.Repository, uid string) (User, error) {
	var u User
	log.Println("the user",uid)

	if _, err := uuid.Parse(uid); err != nil {
	log.Println("the invalid user", uid)

		return u, ErrInvalidID
	}

	stmt := repo.SQ.Select(
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

	q, args, err := stmt.ToSql()
	if err != nil {
		return u, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &u, q, uid); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}

	return u, nil
}

func RetrieveMeByAuthID(ctx context.Context, repo *database.Repository, aid string) (User, error) {
	var u User

	stmt := repo.SQ.Select(
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

	q, args, err := stmt.ToSql()
	if err != nil {
		return u, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &u, q, aid); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}

	return u, nil
}

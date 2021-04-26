package user

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrNotFound = errors.New("user not found")
)

func Create(repo *database.Repository, nu NewUser, now time.Time) (User, error) {
	u := User{
		ID:      uuid.New().String(),
		Auth0ID: nu.Auth0ID,
		Email:   nu.Email,
		Created: now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"users",
	).SetMap(map[string]interface{}{
		"user_id":  u.ID,
		"auth0_id": u.Auth0ID,
		"email":    u.Email,
		"created":  u.Created,
	})

	if _, err := stmt.Exec(); err != nil {
		return u, errors.Wrapf(err, "inserting user: %v")
	}

	return u, nil
}

func Retrieve(repo *database.Repository, aid string) (User, error) {
	var u User

	stmt := repo.SQ.Select(
		"user_id",
		"auth0_id",
		"email",
		"created",
	).From(
		"users",
	).Where(sq.Eq{"auth0_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return u, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.Get(&u, q, aid); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}

	return u, nil
}

func RetrieveByEmail(repo *database.Repository, email string) (User, error) {
	var u User

	stmt := repo.SQ.Select(
		"user_id",
		"auth0_id",
		"email",
		"created",
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

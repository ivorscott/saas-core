package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// The User package shouldn't know anything about http
// While it may identify common know errors, how to respond is left to the handlers
var (
	ErrNotFound       = errors.New("user not found")
	ErrInvalidID      = errors.New("id provided was not a valid UUID")
)

// Create adds a new User
func CreateUser(repo *Repository, nu NewUser, now time.Time) (*User, error) {

	u := User{
		ID:            nu.ID,
		Auth0ID:       nu.Auth0ID,
		Email:         nu.Email,
		EmailVerified: nu.EmailVerified,
		FirstName:     nu.FirstName,
		LastName:      nu.LastName,
		Picture:       nu.Picture,
		Locale:        nu.Locale,
		Created:       now.UTC(),
	}

	stmt := repo.SQ.Insert(
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
		"created":        u.Created,
	}).Suffix("ON CONFLICT DO NOTHING")

	sql,_,_ := stmt.ToSql()
	fmt.Println(sql)

	if _, err := stmt.Exec(); err != nil {
		return nil, errors.Wrapf(err, "inserting user: %v", nu)
	}

	return &u, nil
}

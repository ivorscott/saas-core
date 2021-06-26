package handlers

import (
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/pkg/errors"

	//"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type Users struct {
	repo    database.Storer
	log     *log.Logger
	auth0   auth0.Auther
	origins string
	query   users.Querier
}

// UserQueries defines queries required by user handlers
type UserQueries struct {
	user users.UserQuerier
}

// RetrieveMe retrieves the authenticated user
func (u *User) RetrieveMe(w http.ResponseWriter, r *http.Request) error {
	var us users.User

	uid := u.auth0.UserByID(r.Context())

	if uid == "" {
		return web.NewRequestError(users.ErrNotFound, http.StatusNotFound)
	}
	us, err := u.query.RetrieveMe(r.Context(), u.repo, uid)
	if err != nil {
		switch err {
		case users.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case users.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve authenticated user: %w", err)
		}
	}

	return web.Respond(r.Context(), w, us, http.StatusOK)
}

// Create adds a new user to the internal system and updates the existing Auth0 user
func (u *User) Create(w http.ResponseWriter, r *http.Request) error {
	var nu users.NewUser

	if err := web.Decode(r, &nu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	t, err := u.auth0.GenerateToken()
	if err != nil {
		return err
	}

	var user users.User
	status := http.StatusAccepted

	// does the user already exist?
	user, err = u.query.RetrieveMeByAuthID(r.Context(), u.repo, nu.Auth0ID)
	if err != nil {
		status = http.StatusCreated
		user, err = u.query.Create(r.Context(), u.repo, nu, time.Now())
		if err != nil {
			return err
		}
	}

	if err = u.auth0.UpdateUserAppMetaData(t, nu.Auth0ID, user.ID); err != nil {
		switch err {
		case auth0.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err,"failed to update user app metadata")
		}
	}

	return web.Respond(r.Context(), w, user, status)
}

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
	repo    database.DataStorer
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

	uid := u.auth0.GetUserByID(r.Context())

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

	sub := u.auth0.GetUserBySubject(r.Context())

	// get auth0 management api token
	t, err := u.auth0.GetOrCreateToken()
	if err != nil {
		return err
	}

	// if user already exists update app metadata only
	var us users.User

	us, err = u.query.RetrieveMeByAuthID(r.Context(), u.repo, sub)
	if err == nil {
		if err = u.auth0.UpdateUserAppMetaData(t, sub, us.ID); err != nil {
			return err
		}
		return web.Respond(r.Context(), w, us, http.StatusAccepted)
	}

	if err = web.Decode(r, &nu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	user, err := u.query.Create(r.Context(), u.repo, nu, sub, time.Now())
	if err != nil {
		status = http.StatusCreated
		user, err = u.query.user.Create(r.Context(), u.repo, nu, time.Now())
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	if err = u.auth0.UpdateUserAppMetaData(t, nu.Auth0ID, user.ID); err != nil {
		switch err {
		case auth0.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("failed to update app metadata: %w", err)
		}
	}

	return web.Respond(r.Context(), w, user, status)
}

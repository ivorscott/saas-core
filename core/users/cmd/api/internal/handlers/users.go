package handlers

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"

	"github.com/devpies/devpie-client-core/users/internal/platform/auth0"
	"github.com/devpies/devpie-client-core/users/internal/platform/database"
	"github.com/devpies/devpie-client-core/users/internal/platform/web"
	"github.com/devpies/devpie-client-core/users/internal/users"
)

type Users struct {
	repo    *database.Repository
	log     *log.Logger
	auth0   *auth0.Auth0
	origins string
}

func (u *Users) RetrieveMe(w http.ResponseWriter, r *http.Request) error {
	var us users.User
	var err error

	id := u.auth0.GetUserById(r)

	if id == "" {
		us, err = users.RetrieveMeByAuthID(r.Context(), u.repo, u.auth0.GetUserBySubject(r))
	} else {
		us, err = users.RetrieveMe(r.Context(), u.repo, id)
	}

	if err != nil {
		switch err {
		case users.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case users.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for user %q", id)
		}
	}

	return web.Respond(r.Context(), w, us, http.StatusOK)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) error {
	//sub := u.mauth0.GetUserBySubject(r)

	var t auth0.Token

	var nu users.NewUser
	if err := web.Decode(r, &nu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	//user, err := users.Create(r.Context(), u.repo, nu, sub, time.Now())
	//if err != nil {
	//	return err
	//}
	fmt.Println("Before Retrieval=====================", t)

	// try getting existing auth0 management api token
	t, err := u.auth0.GetToken()
	if err != nil {
		return err
	}

	fmt.Println("After Retrieval=====================", t)

	//if err := u.UpdateUserAppMetaData(t, sub, user.ID); err != nil {
	//	return err
	//}

	return web.Respond(r.Context(), w, nil, http.StatusCreated)
}

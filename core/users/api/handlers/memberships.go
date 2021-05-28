package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/devpies/devpie-client-events/go/events"
)

type Memberships struct {
	repo        *database.Repository
	log         *log.Logger
	auth0       *auth0.Auth0
	nats        *events.Client
}


func (m *Memberships) RetrieveMembers(w http.ResponseWriter, r *http.Request) error {
	uid := m.auth0.GetUserById(r)

	tid := chi.URLParam(r, "tid")

	ms, err := memberships.RetrieveMemberships(r.Context(), m.repo, uid, tid)
	if err != nil {
		switch err {
		case memberships.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case memberships.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for team %q", tid)
		}
	}

	return web.Respond(r.Context(), w, ms, http.StatusOK)
}

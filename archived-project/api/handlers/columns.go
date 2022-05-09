package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/devpies/devpie-client-core/projects/domain/columns"
	"github.com/devpies/devpie-client-core/projects/platform/auth0"
	"github.com/devpies/devpie-client-core/projects/platform/database"
	"github.com/devpies/devpie-client-core/projects/platform/web"
	"github.com/pkg/errors"
)

type Columns struct {
	repo  *database.Repository
	log   *log.Logger
	auth0 *auth0.Auth0
}

func (c *Columns) List(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	list, err := columns.List(r.Context(), c.repo, pid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

func (c *Columns) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	col, err := columns.Retrieve(r.Context(), c.repo, id)
	if err != nil {
		switch err {
		case columns.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case columns.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for columns %q", id)
		}
	}

	return web.Respond(r.Context(), w, col, http.StatusOK)
}

func (c *Columns) Create(w http.ResponseWriter, r *http.Request) error {
	var nc columns.NewColumn
	if err := web.Decode(r, &nc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	col, err := columns.Create(r.Context(), c.repo, nc, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, col, http.StatusCreated)
}

func (c *Columns) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update columns.UpdateColumn
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding column update")
	}

	if err := columns.Update(r.Context(), c.repo, id, update, time.Now()); err != nil {
		switch err {
		case columns.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case columns.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating column %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func (c *Columns) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := columns.Delete(r.Context(), c.repo, id); err != nil {
		switch err {
		case columns.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting column %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

package handlers

import (
	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ivorscott/devpie-client-backend-go/internal/column"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/pkg/errors"
)

// Columns holds the application state needed by the handler methods.
type Columns struct {
	repo *database.Repository
	log  *log.Logger
	auth0 *mid.Auth0
}

// List gets all column
func (c *Columns) List(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	list, err := column.List(r.Context(), c.repo, pid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve a single Column
func (c *Columns) Retrieve(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	id := chi.URLParam(r, "id")

	col, err := column.Retrieve(r.Context(), c.repo, id, pid)
	if err != nil {
		switch err {
		case column.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case column.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for columns %q", id)
		}
	}

	return web.Respond(r.Context(), w, col, http.StatusOK)
}

// Create a new Column
func (c *Columns) Create(w http.ResponseWriter, r *http.Request) error {
	var nc column.NewColumn
	if err := web.Decode(r, &nc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	col, err := column.Create(r.Context(), c.repo, nc, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, col, http.StatusCreated)
}

// Update decodes the body of a request to update an existing column. The ID
// of the column is part of the request URL.
func (c *Columns) Update(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	id := chi.URLParam(r, "id")

	var update column.UpdateColumn
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding column update")
	}

	if err := column.Update(r.Context(), c.repo, pid, id, update); err != nil {
		switch err {
		case column.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case column.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating column %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single column identified by an ID in the request URL.
func (c *Columns) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := column.Delete(r.Context(), c.repo, id); err != nil {
		switch err {
		case column.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting column %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

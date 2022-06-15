package handler

import (
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/devpies/devpie-client-core/projects/domain/columns"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/pkg/errors"
)

type columnService interface {
}

// ColumnHandler handles the column requests.
type ColumnHandler struct {
	logger  *zap.Logger
	service columnService
}

// NewColumnHandler returns a new column handler.
func NewColumnHandler(
	logger *zap.Logger,
	service columnService,
) *ColumnHandler {
	return &ColumnHandler{
		logger:  logger,
		service: service,
	}
}

func (ch *ColumnHandler) List(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	list, err := columns.List(r.Context(), c.repo, pid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

func (ch *ColumnHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
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

func (ch *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) error {
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

func (ch *ColumnHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

func (ch *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) error {
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

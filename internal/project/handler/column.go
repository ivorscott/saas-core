package handler

import (
	"github.com/devpies/saas-core/internal/project"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/devpies/saas-core/pkg/web"
	"github.com/pkg/errors"
)

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
	list, err := ch.service.List(r.Context(), pid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

func (ch *ColumnHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	col, err := ch.service.Retrieve(r.Context(), id)
	if err != nil {
		switch err {
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for columns %q", id)
		}
	}

	return web.Respond(r.Context(), w, col, http.StatusOK)
}

func (ch *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var nc model.NewColumn
	if err := web.Decode(r, &nc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	col, err := ch.service.Create(r.Context(), nc, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, col, http.StatusCreated)
}

func (ch *ColumnHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update model.UpdateColumn
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding column update")
	}

	if err := ch.service.Update(r.Context(), id, update, time.Now()); err != nil {
		switch err {
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating column %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func (ch *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := ch.service.Delete(r.Context(), id); err != nil {
		switch err {
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting column %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

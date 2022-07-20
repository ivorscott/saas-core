package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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

// List handles column list requests.
func (ch *ColumnHandler) List(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	list, err := ch.service.List(r.Context(), pid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve handles column retrieval requests.
func (ch *ColumnHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	col, err := ch.service.Retrieve(r.Context(), id)
	if err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error looking for columns %q: %w", id, err)
		}
	}

	return web.Respond(r.Context(), w, col, http.StatusOK)
}

// Create handles column creation requests.
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

// Update handles column update requests.
func (ch *ColumnHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update model.UpdateColumn
	if err := web.Decode(r, &update); err != nil {
		return fmt.Errorf("error decoding column update: %w", err)
	}

	if _, err := ch.service.Update(r.Context(), id, update, time.Now()); err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error updating column %q: %w", id, err)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// Delete handles column delete requests.
func (ch *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := ch.service.Delete(r.Context(), id); err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error deleting column %q: %w", id, err)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

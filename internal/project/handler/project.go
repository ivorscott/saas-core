package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type projectService interface {
}

// ProjectHandler handles the project requests.
type ProjectHandler struct {
	logger  *zap.Logger
	service projectService
}

// NewProjectHandler returns a new project handler.
func NewProjectHandler(
	logger *zap.Logger,
	service projectService,
) *ProjectHandler {
	return &ProjectHandler{
		logger:  logger,
		service: service,
	}
}

func (ph *ProjectHandler) List(w http.ResponseWriter, r *http.Request) error {
	uid := p.auth0.UserByID(r.Context())
	list, err := projects.List(r.Context(), p.repo, uid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

func (ph *ProjectHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	uid := p.auth0.UserByID(r.Context())
	pid := chi.URLParam(r, "pid")

	opr, err := projects.Retrieve(r.Context(), p.repo, pid, uid)
	if err == nil {
		return web.Respond(r.Context(), w, opr, http.StatusOK)
	}

	spr, err := projects.RetrieveShared(r.Context(), p.repo, pid, uid)
	if err != nil {
		switch err {
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating project %q", pid)
		}
	}

	return web.Respond(r.Context(), w, spr, http.StatusOK)
}

func (ph *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var np projects.NewProject

	uid := p.auth0.UserByID(r.Context())

	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	pr, err := projects.Create(r.Context(), p.repo, np, uid, time.Now())
	if err != nil {
		return err
	}

	e := msg.ProjectCreatedEvent{
		Data: msg.ProjectCreatedEventData{
			ProjectID:   pr.ID,
			Name:        pr.Name,
			Prefix:      pr.Prefix,
			Description: pr.Description,
			TeamID:      pr.TeamID,
			UserID:      pr.UserID,
			Active:      pr.Active,
			Public:      pr.Public,
			ColumnOrder: pr.ColumnOrder,
			UpdatedAt:   pr.UpdatedAt.String(),
			CreatedAt:   pr.CreatedAt.String(),
		},
		Type:     msg.TypeProjectCreated,
		Metadata: msg.Metadata{UserID: uid, TraceID: uuid.New().String()},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	titles := [4]string{"To Do", "In Progress", "Review", "Done"}

	for i, title := range titles {
		nt := columns.NewColumn{
			ProjectID:  pr.ID,
			Title:      title,
			ColumnName: fmt.Sprintf(`column-%d`, i+1),
		}
		_, err := columns.Create(r.Context(), p.repo, nt, time.Now())
		if err != nil {
			return err
		}
	}

	p.nats.Publish(string(msg.TypeProjectCreated), bytes)

	return web.Respond(r.Context(), w, pr, http.StatusCreated)
}

func (ph *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) error {
	var update projects.UpdateProject

	uid := p.auth0.UserByID(r.Context())

	pid := chi.URLParam(r, "pid")

	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding project update")
	}

	up, err := projects.Update(r.Context(), p.repo, pid, uid, update, time.Now())
	if err != nil {
		switch err {
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating project %q", pid)
		}
	}

	e := msg.ProjectUpdatedEvent{
		Type: msg.TypeProjectUpdated,
		Data: msg.ProjectUpdatedEventData{
			Name:        &up.Name,
			Description: &up.Description,
			Active:      &up.Active,
			Public:      &up.Public,
			TeamID:      &up.TeamID,
			ProjectID:   up.ID,
			ColumnOrder: up.ColumnOrder,
			UpdatedAt:   up.UpdatedAt.String(),
		},
		Metadata: msg.Metadata{
			UserID:  uid,
			TraceID: uuid.New().String(),
		},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	p.nats.Publish(string(msg.EventsProjectUpdated), bytes)

	return web.Respond(r.Context(), w, up, http.StatusOK)
}

func (ph *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	uid := p.auth0.UserByID(r.Context())

	if _, err := projects.Retrieve(r.Context(), p.repo, pid, uid); err != nil {
		_, err := projects.RetrieveShared(r.Context(), p.repo, pid, uid)
		if err == nil {
			return web.NewRequestError(err, http.StatusUnauthorized)
		}
	}
	if err := tasks.DeleteAll(r.Context(), p.repo, pid); err != nil {
		return err
	}
	if err := columns.DeleteAll(r.Context(), p.repo, pid); err != nil {
		return err
	}
	if err := projects.Delete(r.Context(), p.repo, pid, uid); err != nil {
		switch err {
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting project %q", pid)
		}
	}

	e := msg.ProjectDeletedEvent{
		Type: msg.TypeProjectDeleted,
		Metadata: msg.Metadata{
			TraceID: uuid.New().String(),
			UserID:  uid,
		},
		Data: msg.ProjectDeletedEventData{
			ProjectID: pid,
		},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	p.nats.Publish(string(msg.EventsProjectDeleted), bytes)

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

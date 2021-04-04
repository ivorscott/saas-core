package handlers

import (
	"fmt"
	"github.com/ivorscott/devpie-client-backend-go/internal/column"
	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"github.com/ivorscott/devpie-client-backend-go/internal/project"
	"github.com/ivorscott/devpie-client-backend-go/internal/task"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/pkg/errors"
)

// Project holds the application state needed by the handler methods.
type Projects struct {
	repo  *database.Repository
	log   *log.Logger
	auth0 *mid.Auth0
}

// List gets all Project
func (p *Projects) List(w http.ResponseWriter, r *http.Request) error {
	uid := p.auth0.GetUserBySubject(r)

	list, err := project.List(r.Context(), p.repo, uid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve a single Project
func (p *Projects) Retrieve(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	uid := p.auth0.GetUserBySubject(r)

	pr, err := project.Retrieve(r.Context(), p.repo, pid, uid)
	if err != nil {
		switch err {
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for projects %q", pid)
		}
	}

	return web.Respond(r.Context(), w, pr, http.StatusOK)
}

// Create a new Project
func (p *Projects) Create(w http.ResponseWriter, r *http.Request) error {
	uid := p.auth0.GetUserBySubject(r)

	var np project.NewProject
	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	pr, err := project.Create(r.Context(), p.repo, np, uid, time.Now())
	if err != nil {
		return err
	}

	// create default columns for project
	titles := [4]string{"To Do", "In Progress", "Review", "Done"}
	for i, title := range titles {
		nt := column.NewColumn{
			ProjectID:  pr.ID,
			Title:      title,
			ColumnName: fmt.Sprintf(`column-%d`, i+1),
		}
		_, err := column.Create(r.Context(), p.repo, nt, time.Now())
		if err != nil {
			return err
		}
	}

	return web.Respond(r.Context(), w, pr, http.StatusCreated)
}

// Update decodes the body of a request to update an existing project. The ID
// of the project is part of the request URL.
func (p *Projects) Update(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	uid := p.auth0.GetUserBySubject(r)

	var update project.UpdateProject
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding project update")
	}

	if err := project.Update(r.Context(), p.repo, pid, update, uid); err != nil {
		switch err {
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating project %q", pid)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single Project identified by an ID in the request URL.
func (p *Projects) Delete(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	uid := p.auth0.GetUserBySubject(r)

	pr, err := project.Retrieve(r.Context(), p.repo, pid, uid)
	if err != nil {
		return err
	}
	if pr != nil {
		if err := task.DeleteAll(r.Context(), p.repo, pid); err != nil {
			return err
		}
		if err := column.DeleteAll(r.Context(), p.repo, pid); err != nil {
			return err
		}
		if err := project.Delete(r.Context(), p.repo, pid, uid); err != nil {
			switch err {
			case project.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return errors.Wrapf(err, "deleting project %q", pid)
			}
		}
	}


	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

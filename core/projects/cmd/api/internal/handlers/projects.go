package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"

	"github.com/devpies/devpie-client-core/projects/internal/columns"
	"github.com/devpies/devpie-client-core/projects/internal/mid"
	"github.com/devpies/devpie-client-core/projects/internal/platform/database"
	"github.com/devpies/devpie-client-core/projects/internal/platform/web"
	"github.com/devpies/devpie-client-core/projects/internal/projects"
	"github.com/devpies/devpie-client-core/projects/internal/tasks"
	"github.com/pkg/errors"
)

type Projects struct {
	repo  *database.Repository
	log   *log.Logger
	auth0 *mid.Auth0
}

func (p *Projects) List(w http.ResponseWriter, r *http.Request) error {
	uid := p.auth0.GetUserBySubject(r)

	list, err := projects.List(r.Context(), p.repo, uid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

func (p *Projects) Retrieve(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	pr, err := projects.Retrieve(r.Context(), p.repo, pid)
	if err != nil {
		switch err {
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for projects %q", pid)
		}
	}

	return web.Respond(r.Context(), w, pr, http.StatusOK)
}

func (p *Projects) Create(w http.ResponseWriter, r *http.Request) error {
	uid := p.auth0.GetUserBySubject(r)

	var np projects.NewProject
	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	pr, err := projects.Create(r.Context(), p.repo, np, uid, time.Now())
	if err != nil {
		return err
	}
	p.log.Print(pr)
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

	return web.Respond(r.Context(), w, pr, http.StatusCreated)
}

func (p *Projects) Update(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	var update projects.UpdateProject
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding project update")
	}

	if _, err := projects.Update(r.Context(), p.repo, pid, update, time.Now()); err != nil {
		switch err {
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating project %q", pid)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

func (p *Projects) Delete(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	uid := p.auth0.GetUserBySubject(r)

	if _, err := projects.Retrieve(r.Context(), p.repo, pid); err != nil {
		return err
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

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

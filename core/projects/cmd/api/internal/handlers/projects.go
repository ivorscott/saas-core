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

type Projects struct {
	repo  *database.Repository
	log   *log.Logger
	auth0 *mid.Auth0
}

func (p *Projects) List(w http.ResponseWriter, r *http.Request) error {
	uid := p.auth0.GetUserBySubject(r)

	list, err := project.List(r.Context(), p.repo, uid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

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
	p.log.Print(pr)
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

func (p *Projects) Update(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	uid := p.auth0.GetUserBySubject(r)

	var update project.UpdateProject
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding project update")
	}

	if _, err := project.Update(r.Context(), p.repo, pid, update, uid); err != nil {
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

func (p *Projects) Delete(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	uid := p.auth0.GetUserBySubject(r)

	if _, err := project.Retrieve(r.Context(), p.repo, pid, uid); err != nil {
		return err
	}
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

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

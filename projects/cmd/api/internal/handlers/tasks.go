package handlers

import (
	"github.com/ivorscott/devpie-client-backend-go/internal/column"
	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/ivorscott/devpie-client-backend-go/internal/task"
	"github.com/pkg/errors"
)

// Tasks holds the application state needed by the handler methods.
type Tasks struct {
	repo  *database.Repository
	log   *log.Logger
	auth0 *mid.Auth0
}

// List gets all task
func (t *Tasks) List(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	list, err := task.List(r.Context(), t.repo, pid)
	if err != nil {
		return err
	}

	if list == nil {
		list = []task.Task{}
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve a single Task
func (t *Tasks) Retrieve(w http.ResponseWriter, r *http.Request) error {

	tid := chi.URLParam(r, "tid")

	ts, err := task.Retrieve(r.Context(), t.repo, tid)
	if err != nil {
		switch err {
		case task.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case task.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for tasks %q", tid)
		}
	}

	return web.Respond(r.Context(), w, ts, http.StatusOK)
}

// Create a new Task
func (t *Tasks) Create(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	cid := chi.URLParam(r, "cid")

	var nt task.NewTask
	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	ts, err := task.Create(r.Context(), t.repo, nt, pid, time.Now())
	if err != nil {
		return err
	}

	c, err := column.Retrieve(r.Context(), t.repo, pid, cid)
	if err != nil {
		return err
	}
	uc := column.UpdateColumn{
		TaskIDS: append(c.TaskIDS, ts.ID),
	}

	if err := column.Update(r.Context(), t.repo, pid, cid, uc); err != nil {
		return err
	}

	return web.Respond(r.Context(), w, ts, http.StatusCreated)
}

// Update decodes the body of a request to update an existing task. The ID
// of the task is part of the request URL.
func (t *Tasks) Update(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	tid := chi.URLParam(r, "tid")

	var ut task.UpdateTask
	if err := web.Decode(r, &ut); err != nil {
		return errors.Wrap(err, "decoding task update")
	}

	if err := task.Update(r.Context(), t.repo, pid, tid, ut); err != nil {
		switch err {
		case task.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case task.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating task %v", ut)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single task identified by an ID in the request URL.
func (t *Tasks) Delete(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	cid := chi.URLParam(r, "cid")
	tid := chi.URLParam(r, "tid")

	c, err := column.Retrieve(r.Context(), t.repo, pid, cid)
	if err != nil {
		return err
	}

	i := SliceIndex(len(c.TaskIDS), func(i int) bool { return c.TaskIDS[i] == tid })

	newTaskIds := append(c.TaskIDS[:i], c.TaskIDS[i+1:]...)
	uc := column.UpdateColumn{TaskIDS: newTaskIds}

	if err := column.Update(r.Context(), t.repo, pid, cid, uc); err != nil {
		return err
	}

	if err := task.Delete(r.Context(), t.repo, pid, tid); err != nil {
		switch err {
		case task.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting task %q", tid)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

func (t *Tasks) Move(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	tid := chi.URLParam(r, "tid")

	var mt task.MoveTask
	if err := web.Decode(r, &mt); err != nil {
		return errors.Wrap(err, "decoding task move")
	}

	cF, err := column.Retrieve(r.Context(), t.repo, pid, mt.From)
	if err != nil {
		return err
	}

	cT, err := column.Retrieve(r.Context(), t.repo, pid, mt.To)
	if err != nil {
		return err
	}

	i := SliceIndex(len(cF.TaskIDS), func(i int) bool { return cF.TaskIDS[i] == tid })

	newFromTaskIds := append(cF.TaskIDS[:i], cF.TaskIDS[i+1:]...)
	foc := column.UpdateColumn{TaskIDS: newFromTaskIds}

	newToTaskIds := append(cT.TaskIDS, tid)
	toc := column.UpdateColumn{TaskIDS: newToTaskIds}

	err = column.Update(r.Context(), t.repo, pid, mt.From, foc)
	err = column.Update(r.Context(), t.repo, pid, mt.To, toc)
	if err != nil {
		switch err {
		case task.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case task.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating column taskIds from:%q, to:%q", mt.From, mt.To)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

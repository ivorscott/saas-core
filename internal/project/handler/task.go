package handler

import (
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/pkg/errors"

	"github.com/devpies/devpie-client-core/projects/domain/columns"
	"github.com/devpies/devpie-client-core/projects/domain/tasks"
	"github.com/devpies/saas-core/pkg/web"
)

type taskService interface{}

// TaskHandler handles the task requests.
type TaskHandler struct {
	logger  *zap.Logger
	service taskService
}

// NewTaskHandler returns a new task handler.
func NewTaskHandler(
	logger *zap.Logger,
	service taskService,
) *TaskHandler {
	return &TaskHandler{
		logger:  logger,
		service: service,
	}
}

func (th *TaskHandler) List(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	list, err := tasks.List(r.Context(), t.repo, pid)
	if err != nil {
		return err
	}

	if list == nil {
		list = []tasks.Task{}
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

func (th *TaskHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {

	tid := chi.URLParam(r, "tid")

	ts, err := tasks.Retrieve(r.Context(), t.repo, tid)
	if err != nil {
		switch err {
		case tasks.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case tasks.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for tasks %q", tid)
		}
	}

	return web.Respond(r.Context(), w, ts, http.StatusOK)
}

func (th *TaskHandler) Create(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")
	cid := chi.URLParam(r, "cid")
	uid := t.auth0.UserByID(r.Context())

	var nt tasks.NewTask
	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	ts, err := tasks.Create(r.Context(), t.repo, nt, pid, uid, time.Now())
	if err != nil {
		return err
	}

	c, err := columns.Retrieve(r.Context(), t.repo, cid)
	if err != nil {
		return err
	}

	ids := append(c.TaskIDS, ts.ID)

	uc := columns.UpdateColumn{
		TaskIDS: &ids,
	}

	if err := columns.Update(r.Context(), t.repo, cid, uc, time.Now()); err != nil {
		return err
	}

	return web.Respond(r.Context(), w, ts, http.StatusCreated)
}

func (th *TaskHandler) Update(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	var ut tasks.UpdateTask
	if err := web.Decode(r, &ut); err != nil {
		return errors.Wrap(err, "decoding task update")
	}

	update, err := tasks.Update(r.Context(), t.repo, tid, ut, time.Now())
	if err != nil {
		switch err {
		case tasks.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case tasks.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating task %v", ut)
		}
	}

	return web.Respond(r.Context(), w, update, http.StatusOK)
}

func (th *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	cid := chi.URLParam(r, "cid")
	tid := chi.URLParam(r, "tid")

	c, err := columns.Retrieve(r.Context(), t.repo, cid)
	if err != nil {
		return err
	}

	i := SliceIndex(len(c.TaskIDS), func(i int) bool { return c.TaskIDS[i] == tid })

	if i >= 0 {
		newTaskIds := append(c.TaskIDS[:i], c.TaskIDS[i+1:]...)
		uc := columns.UpdateColumn{TaskIDS: &newTaskIds}

		if err := columns.Update(r.Context(), t.repo, cid, uc, time.Now()); err != nil {
			return err
		}

		if err := tasks.Delete(r.Context(), t.repo, tid); err != nil {
			switch err {
			case tasks.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return errors.Wrapf(err, "deleting task %q", tid)
			}
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func (th *TaskHandler) Move(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	var mt tasks.MoveTask
	if err := web.Decode(r, &mt); err != nil {
		return errors.Wrap(err, "decoding task move")
	}

	cF, err := columns.Retrieve(r.Context(), t.repo, mt.From)
	if err != nil {
		return err
	}

	cT, err := columns.Retrieve(r.Context(), t.repo, mt.To)
	if err != nil {
		return err
	}

	i := SliceIndex(len(cF.TaskIDS), func(i int) bool { return cF.TaskIDS[i] == tid })

	if i >= 0 {
		newFromTaskIds := append(cF.TaskIDS[:i], cF.TaskIDS[i+1:]...)
		foc := columns.UpdateColumn{TaskIDS: &newFromTaskIds}

		newToTaskIds := append(cT.TaskIDS, tid)
		toc := columns.UpdateColumn{TaskIDS: &newToTaskIds}

		err = columns.Update(r.Context(), t.repo, mt.From, foc, time.Now())
		if err != nil {
			switch err {
			case tasks.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case tasks.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return errors.Wrapf(err, "updating column taskIds from:%q, to:%q", mt.From, mt.To)
			}
		}

		err = columns.Update(r.Context(), t.repo, mt.To, toc, time.Now())
		if err != nil {
			switch err {
			case tasks.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case tasks.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return errors.Wrapf(err, "updating column taskIds from:%q, to:%q", mt.From, mt.To)
			}
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

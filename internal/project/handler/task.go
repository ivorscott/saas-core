package handler

import (
	"github.com/devpies/saas-core/internal/project"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/pkg/errors"

	"github.com/devpies/saas-core/pkg/web"
)

// TaskHandler handles the task requests.
type TaskHandler struct {
	logger        *zap.Logger
	taskService   taskService
	columnService columnService
}

// NewTaskHandler returns a new task handler.
func NewTaskHandler(
	logger *zap.Logger,
	taskService taskService,
	columnService columnService,
) *TaskHandler {
	return &TaskHandler{
		logger:        logger,
		taskService:   taskService,
		columnService: columnService,
	}
}

func (th *TaskHandler) List(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	list, err := th.taskService.List(r.Context(), pid)
	if err != nil {
		return err
	}

	if list == nil {
		list = make([]model.Task, 0)
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

func (th *TaskHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	t, err := th.taskService.Retrieve(r.Context(), tid)
	if err != nil {
		switch err {
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for tasks %q", tid)
		}
	}

	return web.Respond(r.Context(), w, t, http.StatusOK)
}

func (th *TaskHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var (
		uid string
		err error
	)
	pid := chi.URLParam(r, "pid")
	cid := chi.URLParam(r, "cid")

	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}
	uid = values.Metadata.UserID

	var nt model.NewTask
	if err = web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	task, err := th.taskService.Create(r.Context(), nt, pid, uid, time.Now())
	if err != nil {
		return err
	}

	c, err := th.columnService.Retrieve(r.Context(), cid)
	if err != nil {
		return err
	}

	ids := append(c.TaskIDS, task.ID)

	uc := model.UpdateColumn{
		TaskIDS: &ids,
	}

	if err = th.columnService.Update(r.Context(), cid, uc, time.Now()); err != nil {
		return err
	}

	return web.Respond(r.Context(), w, task, http.StatusCreated)
}

func (th *TaskHandler) Update(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	var ut model.UpdateTask
	if err := web.Decode(r, &ut); err != nil {
		return errors.Wrap(err, "decoding task update")
	}

	update, err := th.taskService.Update(r.Context(), tid, ut, time.Now())
	if err != nil {
		switch err {
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating task %v", ut)
		}
	}

	return web.Respond(r.Context(), w, update, http.StatusOK)
}

func (th *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	var err error

	cid := chi.URLParam(r, "cid")
	tid := chi.URLParam(r, "tid")

	c, err := th.columnService.Retrieve(r.Context(), cid)
	if err != nil {
		return err
	}

	i := SliceIndex(len(c.TaskIDS), func(i int) bool { return c.TaskIDS[i] == tid })

	if i >= 0 {
		newTaskIds := append(c.TaskIDS[:i], c.TaskIDS[i+1:]...)
		uc := model.UpdateColumn{TaskIDS: &newTaskIds}

		if err = th.columnService.Update(r.Context(), cid, uc, time.Now()); err != nil {
			return err
		}

		if err = th.taskService.Delete(r.Context(), tid); err != nil {
			switch err {
			case project.ErrInvalidID:
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

	var mt model.MoveTask
	if err := web.Decode(r, &mt); err != nil {
		return errors.Wrap(err, "decoding task move")
	}

	cF, err := th.columnService.Retrieve(r.Context(), mt.From)
	if err != nil {
		return err
	}

	cT, err := th.columnService.Retrieve(r.Context(), mt.To)
	if err != nil {
		return err
	}

	i := SliceIndex(len(cF.TaskIDS), func(i int) bool { return cF.TaskIDS[i] == tid })

	if i >= 0 {
		newFromTaskIds := append(cF.TaskIDS[:i], cF.TaskIDS[i+1:]...)
		foc := model.UpdateColumn{TaskIDS: &newFromTaskIds}

		newToTaskIds := append(cT.TaskIDS, tid)
		toc := model.UpdateColumn{TaskIDS: &newToTaskIds}

		err = th.columnService.Update(r.Context(), mt.From, foc, time.Now())
		if err != nil {
			switch err {
			case project.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case project.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return errors.Wrapf(err, "updating column taskIds from:%q, to:%q", mt.From, mt.To)
			}
		}

		err = th.columnService.Update(r.Context(), mt.To, toc, time.Now())
		if err != nil {
			switch err {
			case project.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case project.ErrInvalidID:
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

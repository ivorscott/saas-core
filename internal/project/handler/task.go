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

// List handles list task requests.
func (th *TaskHandler) List(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	list, err := th.taskService.List(r.Context(), pid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve handles retrieve task requests.
func (th *TaskHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	t, err := th.taskService.Retrieve(r.Context(), tid)
	if err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error looking for tasks %q :%w", tid, err)
		}
	}

	return web.Respond(r.Context(), w, t, http.StatusOK)
}

// Create handles create task requests.
func (th *TaskHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var (
		err error
	)
	pid := chi.URLParam(r, "pid")
	cid := chi.URLParam(r, "cid")

	var nt model.NewTask
	if err = web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	task, err := th.taskService.Create(r.Context(), nt, pid, time.Now())
	if err != nil {
		return err
	}

	c, err := th.columnService.Retrieve(r.Context(), cid)
	if err != nil {
		return err
	}

	ids := append(c.TaskIDS, task.ID)

	uc := model.UpdateColumn{
		TaskIDS: ids,
	}

	if _, err = th.columnService.Update(r.Context(), cid, uc, time.Now()); err != nil {
		return err
	}

	return web.Respond(r.Context(), w, task, http.StatusCreated)
}

// Update handles update task requests.
func (th *TaskHandler) Update(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	var ut model.UpdateTask
	if err := web.Decode(r, &ut); err != nil {
		return fmt.Errorf("error decoding task update : %w", err)
	}

	update, err := th.taskService.Update(r.Context(), tid, ut, time.Now())
	if err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("updating task %v :%w", ut, err)
		}
	}

	return web.Respond(r.Context(), w, update, http.StatusOK)
}

// Delete handles delete task requests.
func (th *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	var err error

	cid := chi.URLParam(r, "cid")
	tid := chi.URLParam(r, "tid")

	c, err := th.columnService.Retrieve(r.Context(), cid)
	if err != nil {
		return err
	}

	i := sliceIndex(len(c.TaskIDS), func(i int) bool { return c.TaskIDS[i] == tid })

	if i >= 0 {
		newTaskIds := append(c.TaskIDS[:i], c.TaskIDS[i+1:]...)
		uc := model.UpdateColumn{TaskIDS: newTaskIds}

		if _, err = th.columnService.Update(r.Context(), cid, uc, time.Now()); err != nil {
			return err
		}

		if err = th.taskService.Delete(r.Context(), tid); err != nil {
			switch err {
			case fail.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return fmt.Errorf("error deleting task %q :%w", tid, err)
			}
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// Move handles move task requests.
func (th *TaskHandler) Move(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	var mt model.MoveTask
	if err := web.Decode(r, &mt); err != nil {
		return fmt.Errorf("decoding task move :%w", err)
	}

	cF, err := th.columnService.Retrieve(r.Context(), mt.From)
	if err != nil {
		return err
	}

	cT, err := th.columnService.Retrieve(r.Context(), mt.To)
	if err != nil {
		return err
	}

	i := sliceIndex(len(cF.TaskIDS), func(i int) bool { return cF.TaskIDS[i] == tid })

	if i >= 0 {
		newFromTaskIds := append(cF.TaskIDS[:i], cF.TaskIDS[i+1:]...)
		foc := model.UpdateColumn{TaskIDS: newFromTaskIds}

		newToTaskIds := append(cT.TaskIDS, tid)
		toc := model.UpdateColumn{TaskIDS: newToTaskIds}

		_, err = th.columnService.Update(r.Context(), mt.From, foc, time.Now())
		if err != nil {
			switch err {
			case fail.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case fail.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return fmt.Errorf("error updating column taskIds from:%q, to:%q : %w", mt.From, mt.To, err)
			}
		}

		_, err = th.columnService.Update(r.Context(), mt.To, toc, time.Now())
		if err != nil {
			switch err {
			case fail.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case fail.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return fmt.Errorf("error updating column taskIds from:%q, to:%q :%w", mt.From, mt.To, err)
			}
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devpies/saas-core/internal/project"
	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// TaskRepository manages data access to project tasks.
type TaskRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewTaskRepository returns a new TaskRepository. The database connection is in the context.
func NewTaskRepository(logger *zap.Logger, pg *db.PostgresDatabase) *TaskRepository {
	return &TaskRepository{
		logger: logger,
		pg:     pg,
	}
}

func (tr *TaskRepository) Retrieve(ctx context.Context, tid string) (model.Task, error) {
	var (
		t   model.Task
		err error
	)

	if _, err = uuid.Parse(tid); err != nil {
		return t, project.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select 
			task_id, key, seq, title, points, content, assigned_to,
			attachments, comments, project_id, updated_at, created_at
		from tasks
		where task_id = ?
	`

	err = conn.QueryRowxContext(ctx, stmt, tid).Scan(&t.ID, &t.Key, &t.Seq, &t.Title, &t.Points, &t.Content, &t.AssignedTo, (*pq.StringArray)(&t.Attachments), (*pq.StringArray)(&t.Comments), &t.ProjectID, &t.UpdatedAt, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, project.ErrNotFound
		}
		return t, err
	}

	return t, nil
}

func (tr *TaskRepository) List(ctx context.Context, pid string) ([]model.Task, error) {
	var (
		t   model.Task
		ts  = make([]model.Task, 0)
		err error
	)

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return ts, project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select 
			task_id, key, seq, title, points, content, assigned_to,
			attachments, comments, project_id, updated_at, created_at
		from tasks
		where project_id = ?
	`

	rows, err := conn.QueryxContext(ctx, stmt, pid)
	if err != nil {
		return nil, errors.Wrap(err, "selecting tasks")
	}
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Key, &t.Seq, &t.Title, &t.Points, &t.Content, &t.AssignedTo, (*pq.StringArray)(&t.Attachments), (*pq.StringArray)(&t.Comments), &t.ProjectID, &t.UpdatedAt, &t.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scanning row into Struct")
		}
		ts = append(ts, t)
	}

	return ts, nil
}

func (tr *TaskRepository) Create(ctx context.Context, nt model.NewTask, pid, uid string, now time.Time) (model.Task, error) {
	var (
		t    model.Task
		last model.Task
		p    model.Project
		err  error
	)

	pr := NewProjectRepository(tr.logger, tr.pg)
	p, err = pr.Retrieve(ctx, pid, uid)
	if err != nil {
		p, err = pr.RetrieveShared(ctx, pid, uid)
		if err != nil {
			return t, err
		}
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, project.ErrConnectionFailed
	}
	defer Close()

	// Get key from last task created in project.
	stmt := `select key from tasks where project_id = ? order by created_at desc limit 1`

	err = conn.QueryRowxContext(ctx, stmt, pid).Scan(&last.Key)
	if err != nil {
		if err != sql.ErrNoRows {
			return t, err
		}
	}
	// Generate sequence number.
	// If no task exists, begin with 1 (e.g., APP-1).
	// Otherwise increment last number.
	var seq = 1
	var lastKeyNumber int
	if last.Key != "" {
		ss := strings.Split(last.Key, "-")
		lastKeyNumber, err = strconv.Atoi(ss[1])
		if err != nil {
			return t, nil
		}
		seq = lastKeyNumber + 1
	}

	k := fmt.Sprintf("%s%d", p.Prefix, seq)

	t = model.Task{
		ID:          uuid.New().String(),
		Key:         k,
		Title:       nt.Title,
		ProjectID:   pid,
		Comments:    make([]string, 0),
		Attachments: make([]string, 0),
	}

	stmt = `
		insert into tasks (task_id, key, title, content, assigned_to, attachments, comments, project_id, updated_at, created_at)
		values (?,?,?,?,?,?,?,?,?,?)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		t.ID,
		t.Key,
		t.Title,
		t.Content,
		t.AssignedTo,
		pq.Array(t.Attachments),
		pq.Array(t.Comments),
		t.ProjectID,
		now.UTC(),
		now.UTC(),
	); err != nil {
		return t, errors.Wrapf(err, "inserting tasks: %v", nt)
	}

	return t, nil
}

func (tr *TaskRepository) Update(ctx context.Context, tid string, update model.UpdateTask, now time.Time) (model.Task, error) {
	var (
		t   model.Task
		err error
	)

	t, err = tr.Retrieve(ctx, tid)
	if err != nil {
		return t, err
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, project.ErrConnectionFailed
	}
	defer Close()

	if update.Title != nil {
		t.Title = *update.Title
	}
	if update.Content != nil {
		t.Content = *update.Content
	}
	if update.AssignedTo != nil {
		t.AssignedTo = *update.AssignedTo
	}
	if update.Attachments != nil {
		t.Attachments = update.Attachments
	}
	if update.Comments != nil {
		t.Comments = update.Comments
	}

	stmt := `
		update tasks
		set
			title = ?,
			content = ?,
			assigned_to = ?,
			comments = ?,
			attachments = ?,
			updated_at = ?
		where task_id = ?
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		t.Title,
		t.Content,
		t.AssignedTo,
		pq.Array(t.Comments),
		pq.Array(t.Attachments),
		now.UTC(),
	); err != nil {
		return t, errors.Wrapf(err, "updating task: %s", tid)
	}

	return t, nil
}

func (tr *TaskRepository) Delete(ctx context.Context, tid string) error {
	var err error

	if _, err = uuid.Parse(tid); err != nil {
		return project.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from tasks where task_id = ?`

	if _, err = conn.ExecContext(ctx, stmt, tid); err != nil {
		return errors.Wrapf(err, "deleting task %s", tid)
	}

	return nil
}

func (tr *TaskRepository) DeleteAll(ctx context.Context, pid string) error {
	var err error

	if _, err = uuid.Parse(pid); err != nil {
		return project.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from tasks where project_id = ?`

	if _, err = conn.ExecContext(ctx, stmt, pid); err != nil {
		return errors.Wrapf(err, "deleting all tasks")
	}

	return nil
}

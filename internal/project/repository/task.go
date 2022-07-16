package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devpies/saas-core/pkg/web"
	"strconv"
	"strings"
	"time"

	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
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

// Retrieve retrieves a specific task from the database.
func (tr *TaskRepository) Retrieve(ctx context.Context, tid string) (model.Task, error) {
	var (
		t   model.Task
		err error
	)

	if _, err = uuid.Parse(tid); err != nil {
		return t, fail.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select 
			task_id, key, seq, title, points, user_id, content, assigned_to,
			attachments, comments, project_id, updated_at, created_at
		from tasks
		where task_id = $1
	`

	err = conn.QueryRowxContext(ctx, stmt, tid).Scan(&t.ID, &t.Key, &t.Seq, &t.Title, &t.Points, &t.UserID, &t.Content, &t.AssignedTo, (*pq.StringArray)(&t.Attachments), (*pq.StringArray)(&t.Comments), &t.ProjectID, &t.UpdatedAt, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, fail.ErrNotFound
		}
		return t, err
	}

	return t, nil
}

// List lists all tasks asscociated to a project.
func (tr *TaskRepository) List(ctx context.Context, pid string) ([]model.Task, error) {
	var (
		t   model.Task
		ts  = make([]model.Task, 0)
		err error
	)

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return ts, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select
			task_id, key, seq, title, points, user_id, content, assigned_to,
			attachments, comments, project_id, updated_at, created_at
		from tasks
		where project_id = $1
	`

	rows, err := conn.QueryxContext(ctx, stmt, pid)
	if err != nil {
		return nil, fmt.Errorf("error selecting tasks: %w", err)
	}

	for rows.Next() {
		err = rows.Scan(
			&t.ID,
			&t.Key,
			&t.Seq,
			&t.Title,
			&t.Points,
			&t.UserID,
			&t.Content,
			&t.AssignedTo,
			(*pq.StringArray)(&t.Attachments),
			(*pq.StringArray)(&t.Comments),
			&t.ProjectID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row into struct: %w", err)
		}
		ts = append(ts, t)
	}

	return ts, nil
}

// Create creates a project task in the database.
func (tr *TaskRepository) Create(ctx context.Context, nt model.NewTask, pid string, now time.Time) (model.Task, error) {
	var (
		t    model.Task
		last model.Task
		p    model.Project
		err  error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return t, web.CtxErr()
	}

	pr := NewProjectRepository(tr.logger, tr.pg)
	p, err = pr.Retrieve(ctx, pid)
	if err != nil {
		if err != nil {
			return t, err
		}
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, fail.ErrConnectionFailed
	}
	defer Close()

	// Get key from last task created in project.
	stmt := `select key from tasks where project_id = $1 order by created_at desc limit 1`

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
		UserID:      values.UserID,
		ProjectID:   pid,
		Comments:    make([]string, 0),
		Attachments: make([]string, 0),
	}

	stmt = `
		insert into tasks (
			task_id, tenant_id, key, title, content, user_id, assigned_to, 
			attachments, comments, project_id, updated_at, created_at
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		t.ID,
		values.TenantID,
		t.Key,
		t.Title,
		t.Content,
		t.UserID,
		t.AssignedTo,
		pq.Array(t.Attachments),
		pq.Array(t.Comments),
		t.ProjectID,
		now.UTC(),
		now.UTC(),
	); err != nil {
		return t, fmt.Errorf("error inserting tasks: %v: %w", nt, err)
	}

	return t, nil
}

// Update updates a specific project task in the database.
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
		return t, fail.ErrConnectionFailed
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
			title = $1,
			content = $2,
			assigned_to = $3,
			comments = $4,
			attachments = $5,
			updated_at = $6
		where task_id = $7
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
		t.ID,
	); err != nil {
		return t, fmt.Errorf("error updating task: %s: %w", tid, err)
	}

	return t, nil
}

// Delete deletes a specific project task from the database.
func (tr *TaskRepository) Delete(ctx context.Context, tid string) error {
	var err error

	if _, err = uuid.Parse(tid); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from tasks where task_id = $1`

	if _, err = conn.ExecContext(ctx, stmt, tid); err != nil {
		return fmt.Errorf("error deleting task %s: %w", tid, err)
	}

	return nil
}

// DeleteAll deletes all project tasks from the database.
func (tr *TaskRepository) DeleteAll(ctx context.Context, pid string) error {
	var err error

	if _, err = uuid.Parse(pid); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from tasks where project_id = $1`

	if _, err = conn.ExecContext(ctx, stmt, pid); err != nil {
		return fmt.Errorf("error deleting all tasks: %w", err)
	}

	return nil
}

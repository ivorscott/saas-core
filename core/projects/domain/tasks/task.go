package tasks

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/projects/domain/projects"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"

	"github.com/devpies/devpie-client-core/projects/platform/database"
)

var (
	ErrNotFound  = errors.New("task not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Retrieve(ctx context.Context, repo *database.Repository, tid string) (Task, error) {
	var t Task

	if _, err := uuid.Parse(tid); err != nil {
		return t, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"task_id",
		"key",
		"seq",
		"title",
		"points",
		"content",
		"assigned_to",
		"attachments",
		"comments",
		"project_id",
		"updated_at",
		"created_at",
	).From(
		"tasks",
	).Where(sq.Eq{"task_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return t, errors.Wrapf(err, "building query: %v", args)
	}

	err = repo.DB.QueryRowContext(ctx, q, tid).Scan(&t.ID, &t.Key, &t.Seq, &t.Title, &t.Points, &t.Content, &t.AssignedTo, (*pq.StringArray)(&t.Attachments), (*pq.StringArray)(&t.Comments), &t.ProjectID, &t.UpdatedAt, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

func List(ctx context.Context, repo *database.Repository, pid string) ([]Task, error) {
	var t Task
	var ts = make([]Task, 0)

	stmt := repo.SQ.Select(
		"task_id",
		"key",
		"seq",
		"title",
		"points",
		"content",
		"assigned_to",
		"attachments",
		"comments",
		"project_id",
		"updated_at",
		"created_at",
	).From("tasks").Where(sq.Eq{"project_id": "?"})
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	rows, err := repo.DB.QueryContext(ctx, q, pid)
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

func Create(ctx context.Context, repo *database.Repository, nt NewTask, pid string, now time.Time) (Task, error) {
	t := Task{
		ID:          uuid.New().String(),
		Title:       nt.Title,
		ProjectID:   pid,
		Comments:    make([]string, 0),
		Attachments: make([]string, 0),
	}

	stmt := repo.SQ.Insert(
		"tasks",
	).SetMap(map[string]interface{}{
		"task_id":     t.ID,
		"key":         t.Key,
		"title":       t.Title,
		"content":     t.Content,
		"assigned_to": t.AssignedTo,
		"attachments": pq.Array(t.Attachments),
		"comments":    pq.Array(t.Comments),
		"project_id":  t.ProjectID,
		"updated_at":  now.UTC(),
		"created_at":  now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return t, errors.Wrapf(err, "inserting tasks: %v", nt)
	}

	// Update Task with Project Prefix and Sequence Number

	p, err := projects.Retrieve(ctx, repo, pid)
	if err != nil {
		return t, err
	}

	task, err := Retrieve(ctx, repo, t.ID)
	if err != nil {
		return t, err
	}

	k := fmt.Sprintf("%s%d", p.Prefix, task.Seq)

	updateStmt := repo.SQ.Update(
		"tasks",
	).SetMap(map[string]interface{}{
		"title":      t.Title,
		"key":        k,
		"updated_at": now.UTC(),
	}).Where(sq.Eq{"task_id": t.ID})

	if _, err := updateStmt.ExecContext(ctx); err != nil {
		return t, errors.Wrapf(err, "updating task %s with key %s", t.ID, k)
	}

	task.Key = k

	return task, nil
}

func Update(ctx context.Context, repo *database.Repository, tid string, update UpdateTask, now time.Time) (Task, error) {
	t, err := Retrieve(ctx, repo, tid)
	if err != nil {
		return t, err
	}

	if update.Title != nil {
		t.Title = *update.Title
	}
	if update.Content != nil {
		t.Content = *update.Content
	}
	if update.Attachments != nil {
		t.Attachments = update.Attachments
	}
	if update.Comments != nil {
		t.Comments = update.Comments
	}

	stmt := repo.SQ.Update(
		"tasks",
	).SetMap(map[string]interface{}{
		"title":       t.Title,
		"content":     t.Content,
		"comments":    pq.Array(t.Comments),
		"attachments": pq.Array(t.Attachments),
		"updated_at":  now.UTC(),
	}).Where(sq.Eq{"task_id": tid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return t, errors.Wrapf(err, "updating task: %s", tid)
	}

	return t, nil
}

func Delete(ctx context.Context, repo *database.Repository, tid string) error {

	if _, err := uuid.Parse(tid); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"tasks",
	).Where(sq.Eq{"task_id": tid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting task %s", tid)
	}

	return nil
}

func DeleteAll(ctx context.Context, repo *database.Repository, pid string) error {

	if _, err := uuid.Parse(pid); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"tasks",
	).Where(sq.Eq{"project_id": pid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting all tasks")
	}

	return nil
}

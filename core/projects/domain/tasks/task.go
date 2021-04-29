package tasks

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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
		"title",
		"content",
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

	if err := repo.DB.GetContext(ctx, &t, q, tid); err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

func List(ctx context.Context, repo *database.Repository, pid string) ([]Task, error) {
	var t = make([]Task, 0)

	stmt := repo.SQ.Select(
		"task_id",
		"title",
		"content",
		"project_id",
		"updated_at",
		"created_at",
	).From("tasks").Where(sq.Eq{"project_id": "?"})
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &t, q, pid); err != nil {
		return nil, errors.Wrap(err, "selecting tasks")
	}

	return t, nil
}

func Create(ctx context.Context, repo *database.Repository, nt NewTask, pid string, now time.Time) (Task, error) {
	t := Task{
		ID:        uuid.New().String(),
		Title:     nt.Title,
		Content:   nt.Content,
		ProjectID: pid,
	}

	stmt := repo.SQ.Insert(
		"tasks",
	).SetMap(map[string]interface{}{
		"task_id":    t.ID,
		"title":      t.Title,
		"content":    t.Content,
		"project_id": t.ProjectID,
		"updated_at":    now.UTC(),
		"created_at":    now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return t, errors.Wrapf(err, "inserting tasks: %v", nt)
	}

	return t, nil
}

func Update(ctx context.Context, repo *database.Repository, tid string, update UpdateTask, now time.Time) error {
	t, err := Retrieve(ctx, repo, tid)
	if err != nil {
		return err
	}

	if update.Title != nil {
		t.Title = *update.Title
	}
	if update.Content != nil {
		t.Content = *update.Content
	}

	stmt := repo.SQ.Update(
		"tasks",
	).SetMap(map[string]interface{}{
		"title":   t.Title,
		"content": t.Content,
		"updated_at": now.UTC(),
	}).Where(sq.Eq{"task_id": tid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "updating task: %s", tid)
	}

	return nil
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

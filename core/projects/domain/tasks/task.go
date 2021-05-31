package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/projects/domain/projects"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"

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

func Create(ctx context.Context, repo *database.Repository, nt NewTask, pid, uid string, now time.Time) (Task, error) {

	var t Task
	var ltk Task
	var p projects.Project

	p, err := projects.Retrieve(ctx, repo, pid, uid)
	if err != nil {
		p, err = projects.RetrieveShared(ctx, repo, pid, uid)
		if err != nil {
			return t, err
		}
	}

	// get key from last task created in project
	stmt1 := repo.SQ.Select(
		"key",
	).From(
		"tasks",
	).Where(sq.Eq{"project_id": pid}).OrderBy(
		"created_at DESC",
	).Limit(1)

	q, args, err := stmt1.ToSql()
	if err != nil {
		return t, errors.Wrapf(err, "building query: %v", args)
	}

	err = repo.DB.QueryRowContext(ctx, q, pid).Scan(&ltk.Key)
	if err != nil {
		if err != sql.ErrNoRows {
			return t, err
		}
	}

	// generate sequence number
	// if no tasks exists than begin with 1 eg., (APP-1)
	// otherwise increment last number
	seq := 1
	if ltk.Key != "" {
		ss := strings.Split(ltk.Key, "-")
		lastKeyNumber, err := strconv.Atoi(ss[1])
		if err != nil {
			return t, nil
		}
		seq = lastKeyNumber + 1
	}

	k := fmt.Sprintf("%s%d", p.Prefix, seq)

	t = Task{
		ID:          uuid.New().String(),
		Key:         k,
		Title:       nt.Title,
		ProjectID:   pid,
		Comments:    make([]string, 0),
		Attachments: make([]string, 0),
	}

	stmt2 := repo.SQ.Insert(
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

	if _, err := stmt2.ExecContext(ctx); err != nil {
		return t, errors.Wrapf(err, "inserting tasks: %v", nt)
	}

	return t, nil
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
	if update.AssignedTo != nil {
		t.AssignedTo = *update.AssignedTo
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
		"assigned_to": t.AssignedTo,
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

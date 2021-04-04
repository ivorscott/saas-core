package column

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

// The Column package shouldn't know anything about http
// While it may identify common know errors, how to respond is left to the handlers
var (
	ErrNotFound  = errors.New("column not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Retrieve(ctx context.Context, repo *database.Repository, pid, cid string) (*Column, error) {
	var c Column

	if _, err := uuid.Parse(cid); err != nil {
		return nil, ErrInvalidID
	}
	if _, err := uuid.Parse(pid); err != nil {
		return nil, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"column_id",
		"project_id",
		"title",
		"column_name",
		"task_ids",
		"created",
	).From(
		"columns",
	).Where(sq.Eq{"column_id": "?", "project_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	err = repo.DB.QueryRowContext(ctx, q, cid, pid).Scan(&c.ID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func List(ctx context.Context, repo *database.Repository, pid string) ([]Column, error) {
	var c Column
	var cs = make([]Column, 0)

	stmt := repo.SQ.Select(
		"column_id",
		"project_id",
		"title",
		"column_name",
		"task_ids",
		"created",
	).From("columns").Where(sq.Eq{"project_id": "?"})
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	rows, err := repo.DB.QueryContext(ctx, q, pid)
	if err != nil {
		return nil, errors.Wrap(err, "selecting columns")
	}
	for rows.Next() {
		err = rows.Scan(&c.ID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.Created)
		if err != nil {
			return nil, errors.Wrap(err, "scanning row into Struct")
		}
		cs = append(cs, c)
	}

	return cs, nil
}

// Create adds a new Column
func Create(ctx context.Context, repo *database.Repository, nc NewColumn, now time.Time) (*Column, error) {
	c := Column{
		ID:         uuid.New().String(),
		Title:      nc.Title,
		ColumnName: nc.ColumnName,
		TaskIDS:    make([]string, 0),
		ProjectID:  nc.ProjectID,
		Created:    now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"columns",
	).SetMap(map[string]interface{}{
		"column_id":   c.ID,
		"title":       c.Title,
		"column_name": c.ColumnName,
		"task_ids":    pq.Array(c.TaskIDS),
		"project_id":  c.ProjectID,
		"created":     now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "inserting column: %v", nc)
	}

	return &c, nil
}

// Update modifies data about a Column. It will error if the specified ID is
// invalid or does not reference an existing Column.
func Update(ctx context.Context, repo *database.Repository, pid, cid string, uc UpdateColumn) error {
	c, err := Retrieve(ctx, repo, pid, cid)
	if err != nil {
		return err
	}

	if uc.Title != nil {
		c.Title = *uc.Title
	}

	c.TaskIDS = uc.TaskIDS

	stmt := repo.SQ.Update(
		"columns",
	).SetMap(map[string]interface{}{
		"title":    c.Title,
		"task_ids": pq.Array(c.TaskIDS),
	}).Where(sq.Eq{"column_id": cid, "project_id": c.ProjectID})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating column")
	}

	return nil
}

// Delete removes the column identified by a given ID.
func Delete(ctx context.Context, repo *database.Repository, cid string) error {
	if _, err := uuid.Parse(cid); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"columns",
	).Where(sq.Eq{"column_id": cid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting column %s", cid)
	}

	return nil
}

// Delete removes all columns identified by pid
func DeleteAll(ctx context.Context, repo *database.Repository, pid string) error {
	if _, err := uuid.Parse(pid); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"columns",
	).Where(sq.Eq{"project_id": pid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting all columns")
	}

	return nil
}

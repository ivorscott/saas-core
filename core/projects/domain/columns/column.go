package columns

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/devpies/devpie-client-core/projects/platform/database"
)

var (
	ErrNotFound  = errors.New("column not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Retrieve(ctx context.Context, repo database.Storer, cid string) (Column, error) {
	var c Column

	if _, err := uuid.Parse(cid); err != nil {
		return c, ErrInvalidID
	}

	stmt := repo.Select(
		"column_id",
		"project_id",
		"title",
		"column_name",
		"task_ids",
		"updated_at",
		"created_at",
	).From(
		"columns",
	).Where(sq.Eq{"column_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return c, errors.Wrapf(err, "building query: %v", args)
	}

	err = repo.QueryRowxContext(ctx, q, cid).Scan(&c.ID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.UpdatedAt, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, ErrNotFound
		}
		return c, err
	}

	return c, nil
}

func List(ctx context.Context, repo database.Storer, pid string) ([]Column, error) {
	var c Column
	var cs = make([]Column, 0)

	stmt := repo.Select(
		"column_id",
		"project_id",
		"title",
		"column_name",
		"task_ids",
		"updated_at",
		"created_at",
	).From("columns").Where(sq.Eq{"project_id": "?"})
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	rows, err := repo.QueryxContext(ctx, q, pid)
	if err != nil {
		return nil, errors.Wrap(err, "selecting columns")
	}
	for rows.Next() {
		err = rows.Scan(&c.ID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.UpdatedAt, &c.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scanning row into Struct")
		}
		cs = append(cs, c)
	}

	return cs, nil
}

func Create(ctx context.Context, repo database.Storer, nc NewColumn, now time.Time) (Column, error) {
	c := Column{
		ID:         uuid.New().String(),
		Title:      nc.Title,
		ColumnName: nc.ColumnName,
		TaskIDS:    make([]string, 0),
		ProjectID:  nc.ProjectID,
		UpdatedAt:  now.UTC(),
		CreatedAt:  now.UTC(),
	}

	stmt := repo.Insert(
		"columns",
	).SetMap(map[string]interface{}{
		"column_id":   c.ID,
		"title":       c.Title,
		"column_name": c.ColumnName,
		"task_ids":    pq.Array(c.TaskIDS),
		"project_id":  c.ProjectID,
		"updated_at":  c.UpdatedAt,
		"created_at":  c.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return c, errors.Wrapf(err, "inserting column: %v", nc)
	}

	return c, nil
}

func Update(ctx context.Context, repo database.Storer, cid string, uc UpdateColumn, now time.Time) error {

	if _, err := uuid.Parse(cid); err != nil {
		return ErrInvalidID
	}

	c, err := Retrieve(ctx, repo, cid)
	if err != nil {
		return err
	}

	if uc.Title != nil {
		c.Title = *uc.Title
	}

	if uc.TaskIDS != nil {
		c.TaskIDS = *uc.TaskIDS
	}

	stmt := repo.Update(
		"columns",
	).SetMap(map[string]interface{}{
		"title":      c.Title,
		"task_ids":   pq.Array(c.TaskIDS),
		"updated_at": now.UTC(),
	}).Where(sq.Eq{"column_id": cid})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating column")
	}

	return nil
}

func Delete(ctx context.Context, repo database.Storer, cid string) error {
	if _, err := uuid.Parse(cid); err != nil {
		return ErrInvalidID
	}
	stmt := repo.Delete(
		"columns",
	).Where(sq.Eq{"column_id": cid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting column %s", cid)
	}

	return nil
}

func DeleteAll(ctx context.Context, repo database.Storer, pid string) error {
	if _, err := uuid.Parse(pid); err != nil {
		return ErrInvalidID
	}

	stmt := repo.Delete(
		"columns",
	).Where(sq.Eq{"project_id": pid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting all columns")
	}

	return nil
}

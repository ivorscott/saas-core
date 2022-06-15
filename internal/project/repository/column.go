package repository

import (
	"context"
	"database/sql"
	"github.com/devpies/saas-core/internal/project"
	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ColumnRepository manages data access to project columns.
type ColumnRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewColumnRepository returns a new ColumnRepository. The database connection is in the context.
func NewColumnRepository(logger *zap.Logger, pg *db.PostgresDatabase) *ColumnRepository {
	return &ColumnRepository{
		logger: logger,
		pg:     pg,
	}
}

func (cr *ColumnRepository) Retrieve(ctx context.Context, cid string) (model.Column, error) {
	var (
		c   model.Column
		err error
	)

	if _, err = uuid.Parse(cid); err != nil {
		return c, project.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return c, project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select column_id, project_id, title, column_name, task_ids, updated_at, created_at
		from columns
		where column_id = ?
	`

	err = conn.QueryRowxContext(ctx, stmt, cid).Scan(&c.ID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.UpdatedAt, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, project.ErrNotFound
		}
		return c, err
	}

	return c, nil
}

func (cr *ColumnRepository) List(ctx context.Context, pid string) ([]model.Column, error) {
	var (
		c   model.Column
		cs  = make([]model.Column, 0)
		err error
	)

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return cs, project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select 
			column_id, project_id, title, column_name, task_ids, updated_at, created_at	
		from columns
		where project_id = ?
	`

	rows, err := conn.QueryxContext(ctx, stmt, pid)
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

func (cr *ColumnRepository) Create(ctx context.Context, nc model.NewColumn, now time.Time) (model.Column, error) {
	var (
		c   model.Column
		err error
	)

	c = model.Column{
		ID:         uuid.New().String(),
		Title:      nc.Title,
		ColumnName: nc.ColumnName,
		TaskIDS:    make([]string, 0),
		ProjectID:  nc.ProjectID,
		UpdatedAt:  now.UTC(),
		CreatedAt:  now.UTC(),
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return c, project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		insert into columns (column_id, title, column_name, task_ids, project_ids, project_id, updated_at, created_at)
		values (?,?,?,?,?,?,?,?)
	`

	if _, err = conn.ExecContext(ctx, stmt, c.ID, c.Title, c.ColumnName, pq.Array(c.TaskIDS), c.ProjectID, c.UpdatedAt, c.CreatedAt); err != nil {
		return c, errors.Wrapf(err, "inserting column: %v", nc)
	}

	return c, nil
}

func (cr *ColumnRepository) Update(ctx context.Context, cid string, uc model.UpdateColumn, now time.Time) error {
	var err error

	if _, err = uuid.Parse(cid); err != nil {
		return project.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	c, err := cr.Retrieve(ctx, cid)
	if err != nil {
		return err
	}

	if uc.Title != nil {
		c.Title = *uc.Title
	}

	if uc.TaskIDS != nil {
		c.TaskIDS = *uc.TaskIDS
	}

	stmt := `
		update columns
		set
			title = ?,
			task_ids = ?,
			updated_at = ?
		where column_id = ?
	`

	_, err = conn.ExecContext(ctx, stmt, c.Title, pq.Array(c.TaskIDS), now.UTC(), cid)
	if err != nil {
		return errors.Wrap(err, "updating column")
	}

	return nil
}

func (cr *ColumnRepository) Delete(ctx context.Context, cid string) error {
	var err error

	if _, err = uuid.Parse(cid); err != nil {
		return project.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from columns where column_id = ?`

	if _, err = conn.ExecContext(ctx, stmt, cid); err != nil {
		return errors.Wrapf(err, "deleting column %s", cid)
	}

	return nil
}

func (cr *ColumnRepository) DeleteAll(ctx context.Context, pid string) error {
	var err error

	if _, err := uuid.Parse(pid); err != nil {
		return project.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from columns where project_id = ?`

	if _, err = conn.ExecContext(ctx, stmt, pid); err != nil {
		return errors.Wrapf(err, "deleting all columns")
	}

	return nil
}

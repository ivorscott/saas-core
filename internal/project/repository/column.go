package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devpies/saas-core/pkg/web"

	"time"

	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
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

// Retrieve retrieves a specific project column from the database.
func (cr *ColumnRepository) Retrieve(ctx context.Context, cid string) (model.Column, error) {
	var (
		c   model.Column
		err error
	)

	if _, err = uuid.Parse(cid); err != nil {
		return c, fail.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return c, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select column_id, project_id, title, column_name, task_ids, updated_at, created_at
		from columns
		where column_id = $1
	`

	err = conn.QueryRowxContext(ctx, stmt, cid).Scan(&c.ID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.UpdatedAt, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, fail.ErrNotFound
		}
		return c, err
	}

	return c, nil
}

// List lists all columns for a project in the database.
func (cr *ColumnRepository) List(ctx context.Context, pid string) ([]model.Column, error) {
	var (
		c   model.Column
		cs  = make([]model.Column, 0)
		err error
	)

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return cs, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select 
			column_id, project_id, title, column_name, task_ids, updated_at, created_at	
		from columns
		where project_id = $1
	`

	rows, err := conn.QueryxContext(ctx, stmt, pid)
	if err != nil {
		return nil, fmt.Errorf("error selecting columns :%w", err)
	}
	for rows.Next() {
		err = rows.Scan(&c.ID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.UpdatedAt, &c.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row into struct :%w", err)
		}
		cs = append(cs, c)
	}

	return cs, nil
}

// Create creates a project column from the database.
func (cr *ColumnRepository) Create(ctx context.Context, nc model.NewColumn, now time.Time) (model.Column, error) {
	var (
		c   model.Column
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return c, web.CtxErr()
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return c, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		insert into columns (
			column_id, tenant_id, title, column_name, task_ids,
			project_id, updated_at, created_at
	 	) values ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		uuid.New().String(),
		values.Metadata.TenantID,
		nc.Title,
		nc.ColumnName,
		pq.Array(make([]string, 0)),
		nc.ProjectID,
		now.UTC(),
		now.UTC(),
	); err != nil {
		return c, fmt.Errorf("error inserting column: %+v :%w", nc, err)
	}

	return c, nil
}

// Update updates a project column from the database.
func (cr *ColumnRepository) Update(ctx context.Context, cid string, uc model.UpdateColumn, now time.Time) error {
	var err error

	if _, err = uuid.Parse(cid); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
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
			title = $1,
			task_ids = $2,
			updated_at = $3
		where column_id = $4
	`

	_, err = conn.ExecContext(ctx, stmt, c.Title, pq.Array(c.TaskIDS), now.UTC(), cid)
	if err != nil {
		return fmt.Errorf("error updating column :%w", err)
	}

	return nil
}

// Delete deletes a project column from the database.
func (cr *ColumnRepository) Delete(ctx context.Context, cid string) error {
	var err error

	if _, err = uuid.Parse(cid); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from columns where column_id = $1`

	if _, err = conn.ExecContext(ctx, stmt, cid); err != nil {
		return fmt.Errorf("error deleting column %s :%w", cid, err)
	}

	return nil
}

// DeleteAll deletes all project columns from the database.
func (cr *ColumnRepository) DeleteAll(ctx context.Context, pid string) error {
	var err error

	if _, err := uuid.Parse(pid); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from columns where project_id = $1`

	if _, err = conn.ExecContext(ctx, stmt, pid); err != nil {
		return fmt.Errorf("error deleting all columns :%w", err)
	}

	return nil
}

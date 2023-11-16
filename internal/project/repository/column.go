package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// ColumnRepository manages data access to project columns.
type ColumnRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewColumnRepository returns a new ColumnRepository.
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
		return c, err
	}
	defer Close()

	stmt := `
		select 
		    column_id, tenant_id, project_id, title, 
		    column_name, task_ids, updated_at, created_at
		from columns
		where column_id = $1
	`

	err = conn.QueryRowxContext(ctx, stmt, cid).Scan(&c.ID, &c.TenantID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.UpdatedAt, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, fail.ErrNotFound
		}
		return c, err
	}

	c.UpdatedAt = c.UpdatedAt.UTC()
	c.CreatedAt = c.CreatedAt.UTC()
	return c, nil
}

// List lists all columns for a project in the database.
func (cr *ColumnRepository) List(ctx context.Context, pid string) ([]model.Column, error) {
	var (
		c   model.Column
		cs  = make([]model.Column, 0)
		err error
	)

	if _, err = uuid.Parse(pid); err != nil {
		return cs, fail.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return cs, err
	}
	defer Close()

	stmt := `
		select 
			column_id, tenant_id, project_id, title, column_name, task_ids, updated_at, created_at	
		from columns
		where project_id = $1
	`

	rows, err := conn.QueryxContext(ctx, stmt, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return cs, nil
		}
		return nil, fmt.Errorf("error selecting columns :%w", err)
	}
	for rows.Next() {
		err = rows.Scan(&c.ID, &c.TenantID, &c.ProjectID, &c.Title, &c.ColumnName, (*pq.StringArray)(&c.TaskIDS), &c.UpdatedAt, &c.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row into struct :%w", err)
		}
		c.UpdatedAt = c.UpdatedAt.UTC()
		c.CreatedAt = c.CreatedAt.UTC()
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
		return c, err
	}
	defer Close()

	c = model.Column{
		ID:         uuid.New().String(),
		TenantID:   values.TenantID,
		Title:      nc.Title,
		ColumnName: nc.ColumnName,
		TaskIDS:    make([]string, 0),
		ProjectID:  nc.ProjectID,
		UpdatedAt:  now.Round(time.Microsecond).UTC(),
		CreatedAt:  now.Round(time.Microsecond).UTC(),
	}

	stmt := `
		insert into columns (
			column_id, tenant_id, title, column_name, task_ids,
			project_id, updated_at, created_at
	 	) values ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		c.ID,
		c.TenantID,
		c.Title,
		c.ColumnName,
		pq.Array(c.TaskIDS),
		c.ProjectID,
		c.UpdatedAt,
		c.CreatedAt,
	); err != nil {
		return c, fmt.Errorf("error inserting column: %+v :%w", nc, err)
	}

	return c, nil
}

// Update updates a project column from the database.
func (cr *ColumnRepository) Update(ctx context.Context, cid string, uc model.UpdateColumn, now time.Time) (model.Column, error) {
	var (
		c   model.Column
		err error
	)

	if _, err = uuid.Parse(cid); err != nil {
		return c, fail.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return c, err
	}
	defer Close()

	c, err = cr.Retrieve(ctx, cid)
	if err != nil {
		return c, err
	}

	if uc.Title != nil {
		c.Title = *uc.Title
	}

	if uc.TaskIDS != nil {
		c.TaskIDS = uc.TaskIDS
	}

	stmt := `
		update columns
		set
			title = $1,
			task_ids = $2,
			updated_at = $3
		where column_id = $4
	`

	_, err = conn.ExecContext(ctx, stmt, c.Title, pq.Array(c.TaskIDS), now.Round(time.Microsecond).UTC(), cid)
	if err != nil {
		return c, fmt.Errorf("error updating column :%w", err)
	}

	return c, nil
}

// Delete deletes a project column from the database.
func (cr *ColumnRepository) Delete(ctx context.Context, cid string) error {
	var err error

	if _, err = uuid.Parse(cid); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := cr.pg.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer Close()

	stmt := `delete from columns where column_id = $1`

	if _, err = conn.ExecContext(ctx, stmt, cid); err != nil {
		return fmt.Errorf("error deleting column %s :%w", cid, err)
	}

	return nil
}

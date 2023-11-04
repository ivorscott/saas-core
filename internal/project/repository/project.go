// Package repository manages the data access layer for handling queries.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/jmoiron/sqlx"
	"regexp"
	"strings"
	"time"

	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// ProjectRepository manages data access to projects.
type ProjectRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
	runTx  func(ctx context.Context, fn func(*sqlx.Tx) error) error
}

// NewProjectRepository returns a new ProjectRepository.
func NewProjectRepository(logger *zap.Logger, pg *db.PostgresDatabase) *ProjectRepository {
	return &ProjectRepository{
		logger: logger,
		pg:     pg,
		runTx:  pg.RunInTransaction,
	}
}

// RunTx runs a function within a transaction context.
func (pr *ProjectRepository) RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	return pr.runTx(ctx, fn)
}

// Retrieve retrieves an owned project from the database.
func (pr *ProjectRepository) Retrieve(ctx context.Context, pid string) (model.Project, error) {
	var (
		p   model.Project
		err error
	)

	if _, err = uuid.Parse(pid); err != nil {
		return p, fail.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, err
	}
	defer Close()

	stmt := `
			select 
				project_id, tenant_id, name, prefix, description,
				user_id, active, "public", column_order, updated_at, created_at
			from projects
			where project_id = $1
		`
	row := conn.QueryRowxContext(ctx, stmt, pid)
	if err = row.Scan(
		&p.ID,
		&p.TenantID,
		&p.Name,
		&p.Prefix,
		&p.Description,
		&p.UserID,
		&p.Active,
		&p.Public,
		(*pq.StringArray)(&p.ColumnOrder),
		&p.UpdatedAt,
		&p.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return p, fail.ErrNotFound
		}
		return p, err
	}

	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()

	return p, nil
}

// List lists a user's projects in the database.
func (pr *ProjectRepository) List(ctx context.Context) ([]model.Project, error) {
	var p model.Project
	var ps = make([]model.Project, 0)

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return ps, err
	}
	defer Close()

	stmt := `select * from projects`

	rows, err := conn.QueryxContext(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("error selecting projects :%w", err)
	}
	for rows.Next() {
		err = rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Prefix, &p.Description, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row into struct :%w", err)
		}

		p.CreatedAt = p.CreatedAt.UTC()
		p.UpdatedAt = p.UpdatedAt.UTC()

		ps = append(ps, p)
	}

	return ps, nil
}

// Create creates a project in the database.
func (pr *ProjectRepository) Create(ctx context.Context, np model.NewProject, now time.Time) (model.Project, error) {
	var (
		p   model.Project
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return p, web.CtxErr()
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, err
	}
	defer Close()

	if _, err = uuid.Parse(values.UserID); err != nil {
		return p, fail.ErrInvalidID
	}

	p = model.Project{
		ID:          uuid.New().String(),
		TenantID:    values.TenantID,
		Name:        np.Name,
		Prefix:      formatPrefix(np.Name),
		Active:      true,
		UserID:      values.UserID,
		ColumnOrder: []string{"column-1", "column-2", "column-3", "column-4"},
		UpdatedAt:   now.UTC(),
		CreatedAt:   now.UTC(),
	}

	stmt := `
			insert into projects (
				project_id, tenant_id, name, prefix,
				description, user_id, column_order, updated_at, created_at
			) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`
	if _, err = conn.ExecContext(
		ctx,
		stmt,
		p.ID,
		p.TenantID,
		p.Name,
		p.Prefix,
		"",
		values.UserID,
		pq.Array(p.ColumnOrder),
		p.UpdatedAt,
		p.CreatedAt,
	); err != nil {
		return p, fmt.Errorf("error inserting project: %+v :%w", p, err)
	}

	return p, nil
}

// Update updates a project in the database.
func (pr *ProjectRepository) Update(ctx context.Context, pid string, update model.UpdateProject, now time.Time) (model.Project, error) {
	var (
		p   model.Project
		err error
	)

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, err
	}
	defer Close()

	p, err = pr.Retrieve(ctx, pid)
	if err != nil {
		return p, err
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Description != nil {
		p.Description = *update.Description
	}
	if update.Active != nil {
		p.Active = *update.Active
	}
	if update.Public != nil {
		p.Public = *update.Public
	}
	if update.ColumnOrder != nil {
		p.ColumnOrder = update.ColumnOrder
	}

	stmt := `
			update projects
			set 
			    name = $1,
			    description = $2,
				active = $3,
				public = $4,
				column_order = $5,
				updated_at = $6
			where project_id = $7
			`

	_, err = conn.ExecContext(
		ctx,
		stmt,
		p.Name,
		p.Description,
		p.Active,
		p.Public,
		pq.Array(p.ColumnOrder),
		now.UTC(),
		pid,
	)
	if err != nil {
		return p, fmt.Errorf("error updating project :%w", err)
	}

	return p, nil
}

// Delete deletes a project, its columns and tasks from the database.
func (pr *ProjectRepository) Delete(ctx context.Context, pid string) error {
	var err error

	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}

	if _, err = uuid.Parse(pid); err != nil {
		return fail.ErrInvalidID
	}

	if _, err = uuid.Parse(values.UserID); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer Close()

	stmt := `delete from tasks where project_id = $1`

	_, err = conn.ExecContext(ctx, stmt, pid)
	if err != nil {
		return fmt.Errorf("error deleting tasks %s :%w", pid, err)
	}

	stmt = `delete from columns where project_id = $1`

	_, err = conn.ExecContext(ctx, stmt, pid)
	if err != nil {
		return fmt.Errorf("error deleting columns %s :%w", pid, err)
	}

	stmt = `delete from projects where project_id = $1 and user_id = $2`

	_, err = conn.ExecContext(ctx, stmt, pid, values.UserID)
	if err != nil {
		return fmt.Errorf("error deleting project %s :%w", pid, err)
	}

	return nil
}

// formatPrefix generates a project prefix.
func formatPrefix(str string) string {
	// Strip spaces and digits.
	var re = regexp.MustCompile(`(\d| )`)
	var substitution = ""

	s := re.ReplaceAllString(str, substitution)
	s = strings.ToUpper(s)

	return s[:3] + "-"
}

package repository

import (
	"context"
	"database/sql"
	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// ProjectRepository manages project data access.
type ProjectRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewProjectRepository returns a new project repository.
func NewProjectRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *ProjectRepository {
	return &ProjectRepository{
		logger: logger,
		pg:     pg,
	}
}

// Create inserts a new project into the database.
func (pr *ProjectRepository) Create(ctx context.Context, p model.ProjectCopy) error {
	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			insert into projects (
				  project_id, tenant_id, name, prefix, description, team_id,
				  user_id, active, public, column_order, updated_at, created_at
			) values (?,?,?,?,?,?,?,?,?,?,?,?)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		p.ID,
		values.Metadata.TenantID,
		p.Name,
		p.Prefix,
		p.Description,
		p.TeamID,
		p.UserID,
		p.Active,
		p.Public,
		pq.Array(p.ColumnOrder),
		p.UpdatedAt,
		p.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

// Retrieve retrieves a single project from the database.
func (pr *ProjectRepository) Retrieve(ctx context.Context, pid string) (model.ProjectCopy, error) {
	var (
		p   model.ProjectCopy
		err error
	)

	if _, err = uuid.Parse(pid); err != nil {
		return p, fail.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			select
			    project_id, tenant_id, name, prefix,description,
			    team_id, user_id, active, public, column_order, updated_at, created_at 
			from projects
			where project_id = ?
	`

	row := conn.QueryRowxContext(ctx, stmt, pid)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fail.ErrNotFound
		}
		return p, err
	}

	return p, nil
}

// Update modifies a project in the database.
func (pr *ProjectRepository) Update(ctx context.Context, pid string, update model.UpdateProjectCopy) error {
	p, err := pr.Retrieve(ctx, pid)
	if err != nil {
		return err
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

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
	if update.TeamID != nil {
		p.TeamID = *update.TeamID
	}
	if update.ColumnOrder != nil {
		p.ColumnOrder = update.ColumnOrder
	}

	stmt := `
			update projects
			set 
			    name = ?,
			    description = ?,
			    active = ?,
			    public = ?,
			    column_order = ?,
			    team_id = ?,
			    updated_at = ?
			where project_id = ?
	`

	_, err = conn.ExecContext(
		ctx,
		stmt,
		p.Name,
		p.Description,
		p.Active,
		p.Public,
		pq.Array(p.ColumnOrder),
		p.TeamID,
		update.UpdatedAt,
		pid,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a project from the database.
func (pr *ProjectRepository) Delete(ctx context.Context, pid string) error {
	var err error

	if _, err = uuid.Parse(pid); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from projects where project_id = ?`

	if _, err = conn.ExecContext(ctx, stmt, pid); err != nil {
		return err
	}

	return nil
}

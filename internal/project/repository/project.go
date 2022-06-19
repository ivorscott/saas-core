package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/jmoiron/sqlx"
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

// NewProjectRepository returns a new ProjectRepository. The database connection is in the context.
func NewProjectRepository(logger *zap.Logger, pg *db.PostgresDatabase) *ProjectRepository {
	return &ProjectRepository{
		logger: logger,
		pg:     pg,
		runTx:  pg.RunInTransaction,
	}
}

func (pr *ProjectRepository) RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	return pr.runTx(ctx, fn)
}

// RetrieveTeamID retrieved the teamID associated with a project from the database.
func (pr *ProjectRepository) RetrieveTeamID(ctx context.Context, pid string) (string, error) {
	var (
		teamID string
		err    error
	)

	if _, err = uuid.Parse(pid); err != nil {
		return "", fail.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return "", fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `select team_id from projects where project_id = $1`
	row := conn.QueryRowxContext(ctx, stmt, pid)
	err = row.Scan(&teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fail.ErrNotFound
		}
		return "", err
	}

	return teamID, nil
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

	values, ok := web.FromContext(ctx)
	if !ok {
		return p, web.CtxErr()
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			select 
				project_id, name, prefix, description, team_id,
				user_id, active, "public", column_order, updated_at, created_at
			from projects
			where project_id = $1 and user_id = $2
		`
	row := conn.QueryRowxContext(ctx, stmt, pid, values.UserID)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fail.ErrNotFound
		}
		return p, err
	}

	return p, nil
}

// RetrieveShared retrieves a shared project from the database.
func (pr *ProjectRepository) RetrieveShared(ctx context.Context, pid string) (model.Project, error) {
	var p model.Project

	if _, err := uuid.Parse(pid); err != nil {
		return p, fail.ErrInvalidID
	}

	values, ok := web.FromContext(ctx)
	if !ok {
		return p, web.CtxErr()
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, fail.ErrConnectionFailed
	}
	defer Close()

	tid, err := pr.RetrieveTeamID(ctx, pid)
	if err != nil {
		return p, err
	}

	membershipRepo := NewMembershipRepository(pr.logger, pr.pg)
	m, err := membershipRepo.Retrieve(ctx, values.UserID, tid)
	if err != nil {
		return p, fail.ErrNotAuthorized
	}

	stmt := `
			select 
				project_id, name, prefix, description, team_id,
				user_id, active, public, column_order, updated_at, created_at
			from projects
			where project_id = $1 and team_id = $2
		`

	row := conn.QueryRowxContext(ctx, stmt, pid, m.TeamID)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fail.ErrNotFound
		}
		return p, err
	}

	return p, nil
}

// List lists a user's projects in the database.
func (pr *ProjectRepository) List(ctx context.Context) ([]model.Project, error) {
	var p model.Project
	var ps = make([]model.Project, 0)

	values, ok := web.FromContext(ctx)
	if !ok {
		return ps, web.CtxErr()
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return ps, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select * from projects
		where team_id in (select team_id from memberships where user_id = $1)
		union 
		select * from projects
	 	where user_id = $1
		group by project_id`

	rows, err := conn.QueryxContext(ctx, stmt, values.UserID)
	if err != nil {
		return nil, fmt.Errorf("error selecting projects :%w", err)
	}
	for rows.Next() {
		err = rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Prefix, &p.Description, &p.UserID, &p.TeamID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row into struct :%w", err)
		}
		ps = append(ps, p)
	}

	pr.logger.Info("TEST!!!!!", zap.Any("list", ps))

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
		return p, fail.ErrConnectionFailed
	}
	defer Close()

	p = model.Project{
		ID:          uuid.New().String(),
		TenantID:    values.TenantID,
		Name:        np.Name,
		Prefix:      fmt.Sprintf("%s-", np.Name[:3]),
		Active:      true,
		UserID:      values.UserID,
		TeamID:      np.TeamID,
		ColumnOrder: []string{"column-1", "column-2", "column-3", "column-4"},
		UpdatedAt:   now.UTC(),
		CreatedAt:   now.UTC(),
	}

	stmt := `
			insert into projects (
				project_id, tenant_id, name, prefix, team_id,
				description, user_id, column_order, updated_at, created_at
			) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`
	if _, err = conn.ExecContext(
		ctx,
		stmt,
		p.ID,
		p.TenantID,
		p.Name,
		p.Prefix,
		np.TeamID,
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
		return p, fail.ErrConnectionFailed
	}
	defer Close()

	p, err = pr.Retrieve(ctx, pid)
	if err != nil {
		p, err = pr.RetrieveShared(ctx, pid)
		if err != nil {
			return p, err
		}
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
	if update.TeamID != nil {
		p.TeamID = *update.TeamID
	}

	stmt := `
			update projects
			set 
			    name = $1,
			    description = $2,
				active = $3,
				public = $4,
				column_order = $5,
				team_id = $6,
				updated_at = $7
			where project_id = $8
			`

	_, err = conn.ExecContext(
		ctx,
		stmt,
		p.Name,
		p.Description,
		p.Active,
		p.Public,
		p.ColumnOrder,
		p.TeamID,
		p.UpdatedAt,
		pid,
	)
	if err != nil {
		return p, fmt.Errorf("error updating project :%w", err)
	}

	return p, nil
}

// UpdateTx updates a project in the database.
func (pr *ProjectRepository) UpdateTx(ctx context.Context, tx *sqlx.Tx, pid string, update model.UpdateProject, now time.Time) (model.Project, error) {
	var (
		p   model.Project
		err error
	)

	p, err = pr.Retrieve(ctx, pid)
	if err != nil {
		p, err = pr.RetrieveShared(ctx, pid)
		if err != nil {
			return p, err
		}
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
	if update.TeamID != nil {
		p.TeamID = *update.TeamID
	}

	stmt := `
			update projects
			set 
			    name = $1,
			    description = $2,
				active = $3,
				public = $4,
				column_order = $5,
				team_id = $6,
				updated_at = $7
			where project_id = $8
			`

	_, err = tx.ExecContext(
		ctx,
		stmt,
		p.Name,
		p.Description,
		p.Active,
		p.Public,
		pq.Array(p.ColumnOrder),
		p.TeamID,
		p.UpdatedAt,
		pid,
	)
	if err != nil {
		return p, fmt.Errorf("error updating project :%w", err)
	}

	return p, nil
}

// Delete deletes a project from the database.
func (pr *ProjectRepository) Delete(ctx context.Context, pid string) error {
	var err error

	if _, err = uuid.Parse(pid); err != nil {
		return fail.ErrInvalidID
	}

	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from projects where project_id = $1 and user_id = $2`

	if _, err = conn.ExecContext(ctx, stmt, pid, values.UserID); err != nil {
		return fmt.Errorf("error deleting project %s :%w", pid, err)
	}

	return nil
}

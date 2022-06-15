package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devpies/saas-core/internal/project"
	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ProjectRepository manages data access to projects.
type ProjectRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewProjectRepository returns a new ProjectRepository. The database connection is in the context.
func NewProjectRepository(logger *zap.Logger, pg *db.PostgresDatabase) *ProjectRepository {
	return &ProjectRepository{
		logger: logger,
		pg:     pg,
	}
}

func (pr *ProjectRepository) RetrieveTeamID(ctx context.Context, pid string) (string, error) {
	var (
		teamID string
		err    error
	)

	if _, err = uuid.Parse(pid); err != nil {
		return "", project.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return "", project.ErrConnectionFailed
	}
	defer Close()

	stmt := `select team_id from projects where project_id = ?`
	row := conn.QueryRowxContext(ctx, stmt, pid)
	err = row.Scan(&teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", project.ErrNotFound
		}
		return "", err
	}

	return teamID, nil
}

func (pr *ProjectRepository) Retrieve(ctx context.Context, pid, uid string) (model.Project, error) {
	var (
		p   model.Project
		err error
	)

	if _, err = uuid.Parse(pid); err != nil {
		return p, project.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			select 
				project_id, name, prefix, description, team_id,
				user_id, active, public, column_order, updated_at, created_at
			from projects
			where project_id = ?, user_id = ?
		`
	row := conn.QueryRowxContext(ctx, stmt, pid, uid)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, project.ErrNotFound
		}
		return p, err
	}

	return p, nil
}

func (pr *ProjectRepository) RetrieveShared(ctx context.Context, pid, uid string) (model.Project, error) {
	var p model.Project

	if _, err := uuid.Parse(pid); err != nil {
		return p, project.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, project.ErrConnectionFailed
	}
	defer Close()

	tid, err := pr.RetrieveTeamID(ctx, pid)
	if err != nil {
		return p, err
	}

	membershipRepo := NewMembershipRepository(pr.logger, pr.pg)
	m, err := membershipRepo.Retrieve(ctx, uid, tid)
	if err != nil {
		return p, project.ErrNotAuthorized
	}

	stmt := `
			select 
				project_id, name, prefix, description, team_id,
				user_id, active, public, column_order, updated_at, created_at
			from projects
			where project_id = ?, team_id = ?
		`

	row := conn.QueryRowxContext(ctx, stmt, pid, m.TeamID)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, project.ErrNotFound
		}
		return p, err
	}

	return p, nil
}

func (pr *ProjectRepository) List(ctx context.Context, uid string) ([]model.Project, error) {
	var p model.Project
	var ps = make([]model.Project, 0)

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return ps, project.ErrConnectionFailed
	}
	defer Close()

	q := `
		select 
    		project_id, name, prefix, description, team_id,
			user_id, active, public, column_order, updated_at, created_at
    	from projects
		where team_id in (select team_id from memberships where user_id = $1)
		union 
			select * from projects where user_id = $1
		group by project_id`

	rows, err := conn.QueryxContext(ctx, q, uid)
	if err != nil {
		return nil, errors.Wrap(err, "selecting projects")
	}
	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.UserID, &p.TeamID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scanning row into Struct")
		}
		ps = append(ps, p)
	}

	return ps, nil
}

func (pr *ProjectRepository) Create(ctx context.Context, np model.NewProject, uid string, now time.Time) (model.Project, error) {
	var (
		p   model.Project
		err error
	)

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, project.ErrConnectionFailed
	}
	defer Close()

	p = model.Project{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Prefix:      fmt.Sprintf("%s-", np.Name[:3]),
		Active:      true,
		UserID:      uid,
		TeamID:      np.TeamID,
		ColumnOrder: []string{"column-1", "column-2", "column-3", "column-4"},
		UpdatedAt:   now.UTC(),
		CreatedAt:   now.UTC(),
	}

	stmt := `
			insert into projects (
				project_id, name, prefix, team_id,
				description, user_id, column_order,updated_at,created_at
			) values (?,?,?,?,?,?,?,?,?)
			`
	if _, err = conn.ExecContext(
		ctx,
		stmt,
		p.ID,
		p.Name,
		p.Prefix,
		p.TeamID,
		"",
		p.UserID,
		pq.Array(p.ColumnOrder),
		p.UpdatedAt,
		p.CreatedAt,
	); err != nil {
		return p, errors.Wrapf(err, "inserting project: %v", p)
	}

	return p, nil
}

func (pr *ProjectRepository) Update(ctx context.Context, pid, uid string, update model.UpdateProject, now time.Time) (model.Project, error) {
	var (
		p   model.Project
		err error
	)

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return p, project.ErrConnectionFailed
	}
	defer Close()

	p, err = pr.Retrieve(ctx, pid, uid)
	if err != nil {
		p, err = pr.RetrieveShared(ctx, pid, uid)
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
		p.ColumnOrder,
		p.TeamID,
		p.UpdatedAt,
		pid,
	)
	if err != nil {
		return p, errors.Wrap(err, "updating project")
	}

	return p, nil
}

func (pr *ProjectRepository) Delete(ctx context.Context, pid, uid string) error {
	var err error

	if _, err = uuid.Parse(pid); err != nil {
		return project.ErrInvalidID
	}

	conn, Close, err := pr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from projects where project_id = ?, user_id = ?`

	if _, err = conn.ExecContext(ctx, stmt, pid, uid); err != nil {
		return errors.Wrapf(err, "deleting project %s", pid)
	}

	return nil
}

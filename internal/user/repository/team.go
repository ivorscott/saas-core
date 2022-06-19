package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TeamRepository manages team data access.
type TeamRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewTeamRepository returns a new team repository.
func NewTeamRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *TeamRepository {
	return &TeamRepository{
		logger: logger,
		pg:     pg,
	}
}

// Create inserts a new team into the database.
func (tr *TeamRepository) Create(ctx context.Context, nt model.NewTeam, now time.Time) (model.Team, error) {
	var (
		t   model.Team
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return t, web.CtxErr()
	}

	if _, err = uuid.Parse(values.UserID); err != nil {
		return t, fail.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, fail.ErrConnectionFailed
	}
	defer Close()

	t = model.Team{
		ID:        uuid.New().String(),
		TenantID:  values.TenantID,
		Name:      nt.Name,
		UserID:    values.UserID,
		UpdatedAt: now.UTC(),
		CreatedAt: now.UTC(),
	}

	stmt := `
		insert into teams (team_id, tenant_id, name, user_id, updated_at, created_at)
		values ($1, $2, $3, $4, $5, $6)
	`
	if _, err = conn.ExecContext(
		ctx,
		stmt,
		t.ID,
		t.TenantID,
		t.Name,
		t.UserID,
		t.UpdatedAt,
		t.CreatedAt,
	); err != nil {
		return t, err
	}

	return t, nil
}

// Retrieve retrieves a single team from the database.
func (tr *TeamRepository) Retrieve(ctx context.Context, tid string) (model.Team, error) {
	var (
		t   model.Team
		err error
	)

	if _, err = uuid.Parse(tid); err != nil {
		return t, fail.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select 
		    team_id, tenant_id, user_id, name, updated_at, created_at
		from teams
		where team_id = $1
	`

	if err = conn.GetContext(ctx, &t, stmt, tid); err != nil {
		if err == sql.ErrNoRows {
			return t, fail.ErrNotFound
		}
		return t, err
	}

	return t, nil
}

// List retrieves a set of teams from the database.
func (tr *TeamRepository) List(ctx context.Context) ([]model.Team, error) {
	var (
		t   model.Team
		ts  = make([]model.Team, 0)
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return ts, web.CtxErr()
	}

	if _, err = uuid.Parse(values.UserID); err != nil {
		return ts, fail.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return ts, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select
			team_id, tenant_id, user_id, name, updated_at, created_at
		from teams
		where team_id in (select team_id from memberships where user_id = $1)
	`

	rows, err := conn.QueryxContext(ctx, stmt, values.UserID)
	if err != nil {
		return ts, err
	}
	for rows.Next() {
		err = rows.StructScan(&t)
		if err != nil {
			return nil, fmt.Errorf("error scanning row into struct :%w", err)
		}
		ts = append(ts, t)
	}

	return ts, nil
}

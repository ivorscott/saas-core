package repository

import (
	"context"
	"database/sql"
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

	if _, err = uuid.Parse(values.Metadata.UserID); err != nil {
		return t, fail.ErrInvalidID
	}

	conn, Close, err := tr.pg.GetConnection(ctx)
	if err != nil {
		return t, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		insert into teams (team_id, tenant_id, name, user_id, updated_at, created_at)
		values ($1, $2, $3, $4, $5, $6)
	`
	if _, err = conn.ExecContext(
		ctx,
		stmt,
		uuid.New().String(),
		values.Metadata.TenantID,
		nt.Name,
		values.Metadata.UserID,
		now.UTC(),
		now.UTC(),
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

	if err = conn.SelectContext(ctx, &t, stmt, tid); err != nil {
		if err == sql.ErrNoRows {
			return t, fail.ErrNotFound
		}
		return t, err
	}

	return t, nil
}

// List retrieves a set of teams from the database.
func (tr *TeamRepository) List(ctx context.Context, uid string) ([]model.Team, error) {
	var (
		ts  []model.Team
		err error
	)

	if _, err = uuid.Parse(uid); err != nil {
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
			where team_id IN (SELECT team_id FROM memberships WHERE user_id = $1)
	`

	if err = conn.SelectContext(ctx, &ts, stmt, uid); err != nil {
		if err == sql.ErrNoRows {
			return ts, fail.ErrNotFound
		}
		return ts, err
	}

	return ts, nil
}

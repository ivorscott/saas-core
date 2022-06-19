package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// MembershipRepository manages membership data access.
type MembershipRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewMembershipRepository returns a new membership repository.
func NewMembershipRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *MembershipRepository {
	return &MembershipRepository{
		logger: logger,
		pg:     pg,
	}
}

// Create inserts a new Membership into the database.
func (mr *MembershipRepository) Create(ctx context.Context, nm model.NewMembership, now time.Time) (model.Membership, error) {
	var (
		m   model.Membership
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return m, web.CtxErr()
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return m, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			insert into memberships (membership_id, tenant_id, user_id, team_id, role, updated_at, created_at)
			values ($1, $2, $3, $4, $5, $6, $7)
	`

	m = model.Membership{
		ID:        uuid.New().String(),
		TenantID:  values.TenantID,
		UserID:    nm.UserID,
		TeamID:    nm.TeamID,
		Role:      nm.Role,
		UpdatedAt: now.UTC(),
		CreatedAt: now.UTC(),
	}
	if _, err = conn.ExecContext(
		ctx,
		stmt,
		m.ID,
		m.TenantID,
		m.UserID,
		m.TeamID,
		m.Role,
		m.UpdatedAt,
		m.CreatedAt,
	); err != nil {
		return m, err
	}

	return m, nil
}

// RetrieveMemberships retrieves a set of memberships from the database.
func (mr *MembershipRepository) RetrieveMemberships(ctx context.Context, tid string) ([]model.MembershipEnhanced, error) {
	var (
		m   model.MembershipEnhanced
		ms  = make([]model.MembershipEnhanced, 0)
		err error
	)

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return ms, fail.ErrConnectionFailed
	}
	defer Close()

	// Verify user has membership first.
	if _, err = mr.RetrieveMembership(ctx, tid); err != nil {
		return ms, err
	}

	stmt := `
		select 
			membership_id, m.tenant_id, user_id, team_id, email,
			first_name, last_name, picture, role, m.updated_at, m.created_at
		from memberships m
		join users using(user_id)
		where team_id = $1 
	`
	rows, err := conn.QueryxContext(ctx, stmt, tid)
	if err != nil {
		if err == sql.ErrNoRows {
			return ms, nil
		}
		return ms, err
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return ms, fmt.Errorf("error decoding struct: %w", err)
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (mr *MembershipRepository) RetrieveMembership(ctx context.Context, tid string) (model.Membership, error) {
	var (
		m   model.Membership
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return m, web.CtxErr()
	}

	if _, err = uuid.Parse(tid); err != nil {
		return m, fail.ErrInvalidID
	}

	if _, err = uuid.Parse(values.UserID); err != nil {
		return m, fail.ErrInvalidID
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return m, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			select membership_id, tenant_id, user_id, team_id, role, updated_at, created_at 
			from memberships
			where team_id = $1 and user_id = $2
	`

	err = conn.QueryRowxContext(ctx, stmt, tid, values.UserID).StructScan(&m)
	if err != nil {
		if err == sql.ErrNoRows {
			return m, nil
		}
		return m, err
	}

	return m, nil
}

// Update modifies a membership in the database.
func (mr *MembershipRepository) Update(ctx context.Context, tid string, update model.UpdateMembership, now time.Time) error {
	var err error

	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}

	m, err := mr.RetrieveMembership(ctx, tid)
	if err != nil {
		return err
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	if update.Role != nil {
		m.Role = *update.Role
	}

	stmt := `update memberships set role = $1, updated_at = $2 where team_id = $3 and user_id = $4`

	_, err = conn.ExecContext(
		ctx,
		stmt,
		m.Role,
		now.UTC(),
		tid,
		values.UserID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a membership from the database.
func (mr *MembershipRepository) Delete(ctx context.Context, tid string) (string, error) {
	var (
		id  string
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return id, web.CtxErr()
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return id, fail.ErrConnectionFailed
	}
	defer Close()

	_, err = mr.RetrieveMembership(ctx, tid)
	if err != nil {
		return id, fail.ErrNotFound
	}

	stmt := `delete from memberships where team_id = $1 and user_id = $2 returning membership_id`

	row := conn.QueryRowxContext(ctx, stmt, tid, values.UserID)

	if err = row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

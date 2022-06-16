package repository

import (
	"context"
	"database/sql"
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
			values (?,?,?,?,?,?,?)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		uuid.New().String(),
		values.Metadata.TenantID,
		nm.UserID,
		nm.TeamID,
		nm.Role,
		now.UTC(),
		now.UTC(),
	); err != nil {
		return m, err
	}

	return m, nil
}

// RetrieveMemberships retrieves a set of memberships from the database.
func (mr *MembershipRepository) RetrieveMemberships(ctx context.Context, uid, tid string) ([]model.MembershipEnhanced, error) {
	var (
		ms  []model.MembershipEnhanced
		err error
	)
	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return ms, fail.ErrConnectionFailed
	}
	defer Close()

	if _, err = mr.RetrieveMembership(ctx, uid, tid); err != nil {
		return ms, err
	}

	stmt := `
			select 
			    membership_id, tenant_id, user_id, team_id, email,
			    first_name, last_name, picture, role, m.updated_at, m.created_at
			from memberships m
			join users using(user_id)
			where team_id = ? 
	`

	if err = conn.SelectContext(ctx, &ms, stmt, tid); err != nil {
		if err == sql.ErrNoRows {
			return nil, fail.ErrNotFound
		}
		return nil, err
	}

	return ms, nil
}

// RetrieveMembership retrieves a single membership from the database.
func (mr *MembershipRepository) RetrieveMembership(ctx context.Context, uid, tid string) (model.Membership, error) {
	var (
		m   model.Membership
		err error
	)

	if _, err = uuid.Parse(tid); err != nil {
		return m, fail.ErrInvalidID
	}

	if _, err = uuid.Parse(uid); err != nil {
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
			where team_id = ? AND user_id = ?
	`

	err = conn.QueryRowxContext(ctx, stmt, tid, uid).StructScan(&m)
	if err != nil {
		if err == sql.ErrNoRows {
			return m, fail.ErrNotFound
		}
		return m, err
	}

	return m, nil
}

// Update modifies a membership in the database.
func (mr *MembershipRepository) Update(ctx context.Context, tid string, update model.UpdateMembership, uid string, now time.Time) error {
	var err error

	m, err := mr.RetrieveMembership(ctx, tid, uid)
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

	stmt := `update memberships set role = ?, updated_at = ? where team_id = ? AND user_id = ?`

	_, err = conn.ExecContext(
		ctx,
		stmt,
		m.Role,
		now.UTC(),
		tid,
		uid,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a membership from the database.
func (mr *MembershipRepository) Delete(ctx context.Context, tid, uid string) (string, error) {
	var (
		id  string
		err error
	)

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return id, fail.ErrConnectionFailed
	}
	defer Close()

	_, err = mr.RetrieveMembership(ctx, tid, uid)
	if err != nil {
		return id, fail.ErrNotFound
	}

	stmt := `delete from memberships where team_id = ? AND user_id = ? returning membership_id`

	row := conn.QueryRowxContext(ctx, stmt, tid, uid)

	if err = row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

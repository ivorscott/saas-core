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

// InviteRepository manages invite data access.
type InviteRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewInviteRepository returns a new invite repository.
func NewInviteRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *InviteRepository {
	return &InviteRepository{
		logger: logger,
		pg:     pg,
	}
}

// Create inserts new invites into the database
func (ir *InviteRepository) Create(ctx context.Context, ni model.NewInvite, now time.Time) (model.Invite, error) {
	var (
		i   model.Invite
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return i, web.CtxErr()
	}

	conn, Close, err := ir.pg.GetConnection(ctx)
	if err != nil {
		return i, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			insert into invites (invite_id, tenant_id, user_id, team_id, read, accepted, expiration, updated_at, created_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		uuid.New().String(),
		values.Metadata.TenantID,
		ni.UserID,
		ni.TeamID,
		false,
		false,
		now.AddDate(0, 0, 5),
		now.UTC(),
		now.UTC(),
	); err != nil {
		return i, err
	}

	return i, nil
}

// RetrieveInvite retrieves a single invite from the database.
func (ir *InviteRepository) RetrieveInvite(ctx context.Context, uid string, iid string) (model.Invite, error) {
	var (
		i   model.Invite
		err error
	)

	if _, err = uuid.Parse(uid); err != nil {
		return i, fail.ErrInvalidID
	}

	if _, err = uuid.Parse(iid); err != nil {
		return i, fail.ErrInvalidID
	}

	conn, Close, err := ir.pg.GetConnection(ctx)
	if err != nil {
		return i, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select 
		    invite_id, tenant_id, user_id, team_id, read, accepted, expiration, updated_at, created_at
		from invites
		where user_id = $1 and invites = $2
	`

	err = conn.QueryRowxContext(ctx, stmt, uid, iid).StructScan(&i)
	if err != nil {
		if err == sql.ErrNoRows {
			return i, fail.ErrNotFound
		}
		return i, err
	}

	return i, nil
}

// RetrieveInvites retrieves a set of invites from the database.
func (ir *InviteRepository) RetrieveInvites(ctx context.Context, uid string) ([]model.Invite, error) {
	var (
		is  []model.Invite
		err error
	)

	if _, err = uuid.Parse(uid); err != nil {
		return is, fail.ErrInvalidID
	}

	conn, Close, err := ir.pg.GetConnection(ctx)
	if err != nil {
		return is, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			select invite_id, tenant_id, user_id, team_id, read, accepted, expiration, updated_at, created_at
			from invites
			where user_id = $1 and expiration > now()
	`

	if err = conn.SelectContext(ctx, &is, stmt, uid); err != nil {
		if err == sql.ErrNoRows {
			return nil, fail.ErrNotFound
		}
		return nil, err
	}

	return is, nil
}

// Update modifies a single invite in the database.
func (ir *InviteRepository) Update(ctx context.Context, update model.UpdateInvite, uid, iid string, now time.Time) (model.Invite, error) {
	var (
		i   model.Invite
		err error
	)

	i, err = ir.RetrieveInvite(ctx, uid, iid)
	if err != nil {
		return i, fail.ErrNotFound
	}

	conn, Close, err := ir.pg.GetConnection(ctx)
	if err != nil {
		return i, fail.ErrConnectionFailed
	}
	defer Close()

	i.Accepted = update.Accepted
	i.UpdatedAt = now.UTC()

	stmt := `update invites set read = true, accepted = $1, updated_at = $2 where user_id = $3 and invite_id = $4`

	_, err = conn.ExecContext(
		ctx,
		stmt,
		i.Accepted,
		i.UpdatedAt,
		uid,
		i.ID,
	)
	if err != nil {
		return i, err
	}

	return i, nil
}

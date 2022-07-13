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
			insert into invites (invite_id, tenant_id, user_id, read, accepted, expiration, updated_at, created_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		uuid.New().String(),
		values.TenantID,
		ni.UserID,
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
func (ir *InviteRepository) RetrieveInvite(ctx context.Context, iid string) (model.Invite, error) {
	var (
		i   model.Invite
		err error
	)
	values, ok := web.FromContext(ctx)
	if !ok {
		return i, web.CtxErr()
	}

	if _, err = uuid.Parse(values.UserID); err != nil {
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
		    invite_id, tenant_id, user_id, read, accepted, expiration, updated_at, created_at
		from invites
		where user_id = $1 and invites = $2
	`

	err = conn.QueryRowxContext(ctx, stmt, values.UserID, iid).StructScan(&i)
	if err != nil {
		if err == sql.ErrNoRows {
			return i, fail.ErrNotFound
		}
		return i, err
	}

	return i, nil
}

// RetrieveInvites retrieves a set of invites from the database.
func (ir *InviteRepository) RetrieveInvites(ctx context.Context) ([]model.Invite, error) {
	var (
		i   model.Invite
		is  = make([]model.Invite, 0)
		err error
	)
	values, ok := web.FromContext(ctx)
	if !ok {
		return is, web.CtxErr()
	}

	if _, err = uuid.Parse(values.UserID); err != nil {
		return is, fail.ErrInvalidID
	}

	conn, Close, err := ir.pg.GetConnection(ctx)
	if err != nil {
		return is, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
			select invite_id, tenant_id, user_id, read, accepted, expiration, updated_at, created_at
			from invites
			where user_id = $1 and expiration > now()
	`

	rows, err := conn.QueryContext(ctx, stmt, values.UserID)
	if err != nil {
		return is, err
	}
	for rows.Next() {
		err = rows.Scan(&i)
		if err != nil {
			return is, fmt.Errorf("error scanning row into struct :%w", err)
		}
		is = append(is, i)
	}

	return is, nil
}

// Update modifies a single invite in the database.
func (ir *InviteRepository) Update(ctx context.Context, update model.UpdateInvite, iid string, now time.Time) (model.Invite, error) {
	var (
		i   model.Invite
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return i, web.CtxErr()
	}

	i, err = ir.RetrieveInvite(ctx, iid)
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
		values.UserID,
		i.ID,
	)
	if err != nil {
		return i, err
	}

	return i, nil
}

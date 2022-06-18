package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MembershipRepository manages data access to team memberships.
type MembershipRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewMembershipRepository returns a new MembershipRepository. The database connection is in the context.
func NewMembershipRepository(logger *zap.Logger, pg *db.PostgresDatabase) *MembershipRepository {
	return &MembershipRepository{
		logger: logger,
		pg:     pg,
	}
}

// Create creates a membership to a team in the database.
func (mr *MembershipRepository) Create(ctx context.Context, nm model.MembershipCopy) error {
	var err error

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		insert into memberships (membership_id, tenant_id, user_id, team_id, role, updated_at, created_at)
		values ($1, $2, $3, $4, $5, $6)
	`

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		nm.ID,
		nm.TenantID,
		nm.UserID,
		nm.TeamID,
		nm.Role,
		nm.UpdatedAt,
		nm.CreatedAt,
	); err != nil {
		return fmt.Errorf("error inserting membership: %w", err)
	}

	return nil
}

// RetrieveByID retrieves a membership by membership id in the database.
func (mr *MembershipRepository) RetrieveByID(ctx context.Context, mid string) (model.MembershipCopy, error) {
	var (
		m   model.MembershipCopy
		err error
	)

	if _, err = uuid.Parse(mid); err != nil {
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
		where membership_id = $1
	`

	if err = conn.SelectContext(ctx, &m, stmt, mid); err != nil {
		if err == sql.ErrNoRows {
			return m, fail.ErrNotFound
		}
		return m, err
	}

	return m, nil
}

// Retrieve retrieves a specific membership in the database.
func (mr *MembershipRepository) Retrieve(ctx context.Context, uid, tid string) (model.MembershipCopy, error) {
	var (
		m   model.MembershipCopy
		err error
	)

	if _, err = uuid.Parse(uid); err != nil {
		return m, fail.ErrInvalidID
	}
	if _, err = uuid.Parse(tid); err != nil {
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
		where user_id = $1 AND team_id = $2
	`

	err = conn.QueryRowxContext(ctx, stmt, uid, tid).StructScan(&m)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return m, fail.ErrNotFound
		}
		return m, err
	}
	return m, nil
}

// Delete deletes a membership in the database.
func (mr *MembershipRepository) Delete(ctx context.Context, mid string) error {
	var err error

	if _, err = uuid.Parse(mid); err != nil {
		return fail.ErrInvalidID
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from memberships where membership_id = $1`

	if _, err = conn.ExecContext(ctx, stmt, mid); err != nil {
		return fmt.Errorf("error deleting membership :%w", err)
	}

	return nil
}

// Update updates a membership in the database.
func (mr *MembershipRepository) Update(ctx context.Context, mid string, update model.UpdateMembershipCopy) error {
	var err error

	if _, err = mr.RetrieveByID(ctx, mid); err != nil {
		return err
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		update memberships
		set 
			role = $1,
			updated_at = $2
		where memberships_id = $3
	`

	if _, err = conn.ExecContext(ctx, stmt, update.Role, update.UpdatedAt, mid); err != nil {
		return fmt.Errorf("error updating membership :%w", err)
	}

	return nil
}

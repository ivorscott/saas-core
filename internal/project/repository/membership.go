package repository

import (
	"context"
	"database/sql"
	"github.com/devpies/saas-core/internal/project"
	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/model"
	"go.uber.org/zap"
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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

func (mr *MembershipRepository) Create(ctx context.Context, nm model.MembershipCopy) error {
	var err error

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		insert into memberships (membership_id, user_id, team_id, role, updated_at, created_at)
		values (?,?,?,?,?,?)
	`

	if _, err := conn.ExecContext(ctx, stmt, nm.ID, nm.UserID, nm.TeamID, nm.Role, nm.UpdatedAt, nm.CreatedAt); err != nil {
		return errors.Wrapf(err, "inserting membership: %v", err)
	}

	return nil
}

func (mr *MembershipRepository) RetrieveById(ctx context.Context, mid string) (model.MembershipCopy, error) {
	var (
		m   model.MembershipCopy
		err error
	)

	if _, err = uuid.Parse(mid); err != nil {
		return m, project.ErrInvalidID
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return m, project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select membership_id, user_id, team_id, role, updated_at, created_at
		from memberships
		where membership_id = ?
	`

	if err = conn.SelectContext(ctx, &m, stmt, mid); err != nil {
		if err == sql.ErrNoRows {
			return m, project.ErrNotFound
		}
		return m, err
	}

	return m, nil
}

func (mr *MembershipRepository) Retrieve(ctx context.Context, uid, tid string) (model.MembershipCopy, error) {
	var (
		m   model.MembershipCopy
		err error
	)

	if _, err = uuid.Parse(uid); err != nil {
		return m, project.ErrInvalidID
	}
	if _, err = uuid.Parse(tid); err != nil {
		return m, project.ErrInvalidID
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return m, project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		select membership_id, user_id, team_id, role, updated_at, created_at
		from memberships
		where user_id = ? AND team_id = ?
	`

	err = conn.QueryRowxContext(ctx, stmt, uid, tid).StructScan(&m)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return m, project.ErrNotFound
		}
		return m, err
	}
	return m, nil
}

func (mr *MembershipRepository) Update(ctx context.Context, mid string, update model.UpdateMembershipCopy) error {
	var err error

	if _, err = mr.RetrieveById(ctx, mid); err != nil {
		return err
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	stmt := `
		update memberships
		set 
			role = ?,
			updated_at = ?
		where memberships_id = ?
	`

	if _, err = conn.ExecContext(ctx, stmt, update.Role, update.UpdatedAt, mid); err != nil {
		return errors.Wrap(err, "updating membership")
	}

	return nil
}

func (mr *MembershipRepository) Delete(ctx context.Context, mid string) error {
	var err error

	if _, err = uuid.Parse(mid); err != nil {
		return project.ErrInvalidID
	}

	conn, Close, err := mr.pg.GetConnection(ctx)
	if err != nil {
		return project.ErrConnectionFailed
	}
	defer Close()

	stmt := `delete from memberships where membership_id = ?`

	if _, err = conn.ExecContext(ctx, stmt, mid); err != nil {
		return errors.Wrap(err, "deleting membership")
	}

	return nil
}

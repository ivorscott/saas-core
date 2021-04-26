package memberships

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/projects/internal/platform/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrNotFound  = errors.New("membership not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Create(ctx context.Context, repo *database.Repository, nm MembershipCopy) error {
	stmt := repo.SQ.Insert(
		"memberships",
	).SetMap(map[string]interface{}{
		"membership_id": nm.ID,
		"user_id":       nm.UserID,
		"team_id":       nm.TeamID,
		"role":          nm.Role,
		"updated_at":    nm.UpdatedAt,
		"created_at":    nm.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "inserting membership: %v", err)
	}

	return nil
}

func Retrieve(ctx context.Context, repo *database.Repository, mid string) (MembershipCopy, error) {
	var m MembershipCopy

	if _, err := uuid.Parse(mid); err != nil {
		return m, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"membership_id",
		"user_id",
		"team_id",
		"role",
		"updated_at",
		"created_at",
	).From(
		"memberships",
	).Where(sq.Eq{"membership_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return m, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &m, q, mid); err != nil {
		if err == sql.ErrNoRows {
			return m, ErrNotFound
		}
		return m, err
	}

	return m, nil
}

func Update(ctx context.Context, repo *database.Repository, mid string, update UpdateMembershipCopy) error {
	if _, err := Retrieve(ctx, repo, mid); err != nil {
		return err
	}

	stmt := repo.SQ.Update(
		"memberships",
	).SetMap(map[string]interface{}{
		"role":       update.Role,
		"updated_at": update.UpdatedAt,
	}).Where(sq.Eq{"membership_id": mid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrap(err, "updating membership")
	}

	return nil
}

func Delete(ctx context.Context, repo *database.Repository, mid string) error {
	if _, err := uuid.Parse(mid); err != nil {
		return ErrInvalidID
	}
	stmt := repo.SQ.Delete(
		"memberships",
	).Where(sq.Eq{"membership_id": mid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrap(err, "deleting membership")
	}

	return nil
}

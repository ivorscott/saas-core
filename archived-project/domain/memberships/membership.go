package memberships

import (
	"context"
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/devpies/devpie-client-core/projects/platform/database"
)

var (
	ErrNotFound  = errors.New("membership not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Create(ctx context.Context, repo database.Storer, nm MembershipCopy) error {
	stmt := repo.Insert(
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

func RetrieveById(ctx context.Context, repo database.Storer, mid string) (MembershipCopy, error) {
	var m MembershipCopy

	if _, err := uuid.Parse(mid); err != nil {
		return m, ErrInvalidID
	}

	stmt := repo.Select(
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

	if err := repo.SelectContext(ctx, &m, q, mid); err != nil {
		if err == sql.ErrNoRows {
			return m, ErrNotFound
		}
		return m, err
	}

	return m, nil
}

func Retrieve(ctx context.Context, repo database.Storer, uid, tid string) (MembershipCopy, error) {
	var m MembershipCopy

	if _, err := uuid.Parse(uid); err != nil {
		return m, ErrInvalidID
	}
	if _, err := uuid.Parse(tid); err != nil {
		return m, ErrInvalidID
	}

	stmt := repo.Select(
		"membership_id",
		"user_id",
		"team_id",
		"role",
		"updated_at",
		"created_at",
	).From(
		"memberships",
	).Where("user_id = ? AND team_id = ?")

	q, args, err := stmt.ToSql()
	if err != nil {
		return m, errors.Wrapf(err, "building query: %v", args)
	}

	err = repo.QueryRowxContext(ctx, q, uid, tid).StructScan(&m)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return m, ErrNotFound
		}
		return m, err
	}
	return m, nil
}

func Update(ctx context.Context, repo database.Storer, mid string, update UpdateMembershipCopy) error {
	if _, err := RetrieveById(ctx, repo, mid); err != nil {
		return err
	}

	stmt := repo.Update(
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

func Delete(ctx context.Context, repo database.Storer, mid string) error {
	if _, err := uuid.Parse(mid); err != nil {
		return ErrInvalidID
	}
	stmt := repo.Delete(
		"memberships",
	).Where(sq.Eq{"membership_id": mid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrap(err, "deleting membership")
	}

	return nil
}

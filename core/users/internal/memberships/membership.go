package memberships

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-core/users/internal/platform/database"
	"github.com/ivorscott/devpie-client-core/users/internal/teams"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrNotFound = errors.New("membership not found")
)

func Create(ctx context.Context, repo *database.Repository, nm NewMembership, now time.Time) (Membership, error) {
	m := Membership{
		ID:      uuid.New().String(),
		UserID:  nm.UserID,
		TeamID:  nm.TeamID,
		Role:    nm.Role,
	}

	stmt := repo.SQ.Insert(
		"memberships",
	).SetMap(map[string]interface{}{
		"membership_id": m.ID,
		"user_id":       m.UserID,
		"team_id":       m.TeamID,
		"role":          m.Role,
		"updated_at":    now.UTC(),
		"created_at":    now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return m, errors.Wrapf(err, "inserting membership: %v", err)
	}

	return m, nil
}

func RetrieveMemberships(ctx context.Context, repo *database.Repository, uid string) ([]Membership, error) {
	var m []Membership

	stmt := repo.SQ.Select(
		"membership_id",
		"user_id",
		"team_id",
		"role",
		"updated_at",
		"created_at",
	).From(
		"memberships",
	).Where(sq.Eq{"user_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &m, q, uid); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return m, nil
}

func RetrieveMembership(ctx context.Context, repo *database.Repository, tid, uid string) (Membership, error) {
	var m Membership

	if _, err := uuid.Parse(tid); err != nil {
		return m, teams.ErrInvalidID
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
	).Where(sq.Eq{"team_id": "?", "user_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return m, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &m, q, tid, uid); err != nil {
		if err == sql.ErrNoRows {
			return m, ErrNotFound
		}
		return m, err
	}

	return m, nil
}

func Update(ctx context.Context, repo *database.Repository, tid string, update UpdateMembership, uid string, now time.Time) error {
	m, err := RetrieveMembership(ctx, repo, tid, uid)
	if err != nil {
		return err
	}

	if update.Role != nil {
		m.Role = *update.Role
	}

	stmt := repo.SQ.Update(
		"memberships",
	).SetMap(map[string]interface{}{
		"role":       m.Role,
		"updated_at": now.UTC(),
	}).Where(sq.Eq{"team_id": tid, "user_id": uid})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating membership")
	}

	return nil
}

package invite

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/pkg/errors"
	"time"
)

// Create adds a new Team
func CreateInvites(ctx context.Context, repo *database.Repository, nm NewMembership, now time.Time) (Membership, error) {
	m := Membership{
		ID:      uuid.New().String(),
		UserID:  nm.UserID,
		TeamID:  nm.TeamID,
		Role:    nm.Role,
		Created: now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"membership",
	).SetMap(map[string]interface{}{
		"membership_id": m.ID,
		"user_id":   m.UserID,
		"team_id":   m.TeamID,
		"role":      m.Role,
		"created":   now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return m, errors.Wrapf(err, "inserting membership: %v", err)
	}

	return m, nil
}

func RetrieveInvites(ctx context.Context, repo *database.Repository, uid string) ([]Membership, error) {
	var m []Membership

	stmt := repo.SQ.Select(
		"membership_id",
		"user_id",
		"team_id",
		"role",
		"created",
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

package teams

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrNotFound  = errors.New("team not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Create(ctx context.Context, repo *database.Repository, nt NewTeam, uid string, now time.Time) (Team, error) {
	t := Team{
		ID:      uuid.New().String(),
		Name:    nt.Name,
		UserId:  uid,
	}

	stmt := repo.SQ.Insert(
		"teams",
	).SetMap(map[string]interface{}{
		"team_id": t.ID,
		"name":    t.Name,
		"user_id": t.UserId,
		"updated_at": now.UTC(),
		"created_at": now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return t, errors.Wrap(err, "inserting team")
	}

	return t, nil
}

func Retrieve(ctx context.Context, repo *database.Repository, tid string) (Team, error) {
	var t Team

	if _, err := uuid.Parse(tid); err != nil {
		return t, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"team_id",
		"user_id",
		"name",
		"updated_at",
		"created_at",
	).From(
		"teams",
	).Where(sq.Eq{"team_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return t, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &t, q, tid); err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

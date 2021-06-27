package teams

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/google/uuid"
)

// Error codes returned by failures to handle teams.
var (
	ErrNotFound  = errors.New("team not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

type TeamQuerier interface {
	Create(ctx context.Context, repo database.Storer, nt NewTeam, uid string, now time.Time) (Team, error)
	Retrieve(ctx context.Context, repo database.Storer, tid string) (Team, error)
	List(ctx context.Context, repo database.Storer, uid string) ([]Team, error)
}

type Queries struct{}

func (q *Queries) Create(ctx context.Context, repo database.Storer, nt NewTeam, uid string, now time.Time) (Team, error) {
	var t Team

	if _, err := uuid.Parse(uid); err != nil {
		return t, ErrInvalidID
	}

	t = Team{
		ID:        uuid.New().String(),
		Name:      nt.Name,
		UserID:    uid,
		UpdatedAt: now.UTC(),
		CreatedAt: now.UTC(),
	}

	stmt := repo.Insert(
		"teams",
	).SetMap(map[string]interface{}{
		"team_id":    t.ID,
		"name":       t.Name,
		"user_id":    t.UserID,
		"updated_at": t.UpdatedAt,
		"created_at": t.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return t, err
	}

	return t, nil
}

func (q *Queries) Retrieve(ctx context.Context, repo database.Storer, tid string) (Team, error) {
	var t Team

	if _, err := uuid.Parse(tid); err != nil {
		return t, ErrInvalidID
	}

	stmt := repo.Select(
		"team_id",
		"user_id",
		"name",
		"updated_at",
		"created_at",
	).From(
		"teams",
	).Where(sq.Eq{"team_id": "?"})

	query, args, err := stmt.ToSql()
	if err != nil {
		return t, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	if err := repo.GetContext(ctx, &t, query, tid); err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

func (q *Queries) List(ctx context.Context, repo database.Storer, uid string) ([]Team, error) {
	var ts []Team

	if _, err := uuid.Parse(uid); err != nil {
		return ts, ErrInvalidID
	}

	stmt := repo.Select(
		"team_id",
		"user_id",
		"name",
		"updated_at",
		"created_at",
	).From(
		"teams",
	).Where("team_id IN (SELECT team_id FROM memberships WHERE user_id = ?)")

	query, args, err := stmt.ToSql()
	if err != nil {
		return ts, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	if err := repo.SelectContext(ctx, &ts, query, uid); err != nil {
		if err == sql.ErrNoRows {
			return ts, ErrNotFound
		}
		return ts, err
	}

	return ts, nil
}

package invites

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
	ErrNotFound = errors.New("invite not found")
)

func Create(ctx context.Context, repo *database.Repository, ni NewInvite, now time.Time) (Invite, error) {
	i := Invite{
		ID:         uuid.New().String(),
		UserID:     ni.UserID,
		TeamID:     ni.TeamID,
		Read:       false,
		Accepted:   false,
		Expiration: now.AddDate(0, 0, 5),
	}

	stmt := repo.SQ.Insert(
		"invites",
	).SetMap(map[string]interface{}{
		"invite_id":  i.ID,
		"user_id":    i.UserID,
		"team_id":    i.TeamID,
		"read":       i.Read,
		"accepted":   i.Accepted,
		"expiration": i.Expiration,
		"updated_at": now.UTC(),
		"created_at": now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return i, errors.Wrapf(err, "inserting invite: %v", err)
	}

	return i, nil
}

func RetrieveInvite(ctx context.Context, repo *database.Repository, uid, iid string) (Invite, error) {
	var i Invite

	stmt := repo.SQ.Select(
		"invite_id",
		"user_id",
		"team_id",
		"read",
		"activated",
		"expiration",
		"updated_at",
		"created_at",
	).From(
		"invites",
	).Where(sq.Eq{"user_id": "?", "invite_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return i, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &i, q, uid, iid); err != nil {
		if err == sql.ErrNoRows {
			return i, ErrNotFound
		}
		return i, err
	}

	return i, nil
}

func RetrieveInvites(ctx context.Context, repo *database.Repository, uid string, now time.Time) ([]Invite, error) {
	var is []Invite
	// TODO: return invites that have not expired
	stmt := repo.SQ.Select(
		"invite_id",
		"user_id",
		"team_id",
		"read",
		"activated",
		"expiration",
		"updated_at",
		"created_at",
	).From(
		"invites",
	).Where(sq.Eq{"user_id": "?"}).Where("expiration > ?")

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &is, q, uid, now.UTC()); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return is, nil
}

func Update(ctx context.Context, repo *database.Repository, update UpdateInvite, uid, iid string, now time.Time) error {
	i, err := RetrieveInvite(ctx, repo, uid, iid)
	if err != nil {
		return err
	}

	if update.Read != nil {
		i.Read = *update.Read
	}
	if update.Accepted != nil {
		i.Accepted = *update.Accepted
	}
	if update.Expiration != nil {
		i.Expiration = *update.Expiration
	}

	stmt := repo.SQ.Update(
		"invites",
	).SetMap(map[string]interface{}{
		"read":       i.Read,
		"accepted":   i.Accepted,
		"expiration": i.Expiration,
		"updated_at": now.UTC,
	}).Where(sq.Eq{"user_id": uid, "invite_id": iid})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating invite")
	}

	return nil
}

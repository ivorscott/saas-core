package invites

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

// Error codes returned by failures to handle invites.
var (
	ErrNotFound  = errors.New("invite not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

// InviteQuerier describes the behavior required for executing Invite related queries
type InviteQuerier interface {
	Create(ctx context.Context, repo database.Storer, ni NewInvite, now time.Time) (Invite, error)
	RetrieveInvite(ctx context.Context, repo database.Storer, uid string, iid string) (Invite, error)
	RetrieveInvites(ctx context.Context, repo database.Storer, uid string) ([]Invite, error)
	Update(ctx context.Context, repo database.Storer, update UpdateInvite, uid, iid string, now time.Time) (Invite, error)
}

// Queries defines method implementations for interacting with the invites table
type Queries struct{}

// Create inserts new invites into the database
func (q *Queries) Create(ctx context.Context, repo database.Storer, ni NewInvite, now time.Time) (Invite, error) {
	i := Invite{
		ID:         uuid.New().String(),
		UserID:     ni.UserID,
		TeamID:     ni.TeamID,
		Read:       false,
		Accepted:   false,
		Expiration: now.AddDate(0, 0, 5),
		UpdatedAt:  now.UTC(),
		CreatedAt:  now.UTC(),
	}

	stmt := repo.Insert(
		"invites",
	).SetMap(map[string]interface{}{
		"invite_id":  i.ID,
		"user_id":    i.UserID,
		"team_id":    i.TeamID,
		"read":       i.Read,
		"accepted":   i.Accepted,
		"expiration": i.Expiration,
		"updated_at": i.UpdatedAt,
		"created_at": i.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return i, err
	}

	return i, nil
}

// RetrieveInvite retrieves a single invite from the database
func (q *Queries) RetrieveInvite(ctx context.Context, repo database.Storer, uid string, iid string) (Invite, error) {
	var i Invite

	if _, err := uuid.Parse(uid); err != nil {
		return i, ErrInvalidID
	}

	if _, err := uuid.Parse(iid); err != nil {
		return i, ErrInvalidID
	}

	stmt := repo.Select(
		"invite_id",
		"user_id",
		"team_id",
		"read",
		"accepted",
		"expiration",
		"updated_at",
		"created_at",
	).From(
		"invites",
	).Where("user_id = ? AND invite_id = ?")

	query, args, err := stmt.ToSql()
	if err != nil {
		return i, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	err = repo.QueryRowxContext(ctx, query, uid, iid).StructScan(&i)
	if err != nil {
		if err == sql.ErrNoRows {
			return i, ErrNotFound
		}
		return i, err
	}

	return i, nil
}

// RetrieveInvites retrieves a set of invites from the database
func (q *Queries) RetrieveInvites(ctx context.Context, repo database.Storer, uid string) ([]Invite, error) {
	var is []Invite

	if _, err := uuid.Parse(uid); err != nil {
		return is, ErrInvalidID
	}

	stmt := repo.Select(
		"invite_id",
		"user_id",
		"team_id",
		"read",
		"accepted",
		"expiration",
		"updated_at",
		"created_at",
	).From(
		"invites",
	).Where(sq.Eq{"user_id": "?"}).Where("expiration > NOW()")

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	if err := repo.SelectContext(ctx, &is, query, uid); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return is, nil
}

// Update modifies a single invite in the database
func (q *Queries) Update(ctx context.Context, repo database.Storer, update UpdateInvite, uid, iid string, now time.Time) (Invite, error) {
	i, err := q.RetrieveInvite(ctx, repo, uid, iid)
	if err != nil {
		return i, err
	}

	i.Accepted = update.Accepted
	i.UpdatedAt = now.UTC()

	stmt := repo.Update(
		"invites",
	).SetMap(map[string]interface{}{
		"read":       true,
		"accepted":   i.Accepted,
		"updated_at": i.UpdatedAt,
	}).Where(sq.Eq{"user_id": uid, "invite_id": i.ID})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return i, err
	}

	return i, nil
}

package team

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/pkg/errors"
	"time"
)

// Create adds a new Team
func CreateMember(ctx context.Context, repo *database.Repository, nm NewMember, now time.Time) (Member, error) {
	m := Member{
		ID:             uuid.New().String(),
		UserID:         nm.UserID,
		TeamID:         nm.TeamID,
		IsLeader:       nm.IsLeader,
		InviteSent:     nm.InviteSent,
		InviteAccepted: nm.InviteAccepted,
		Created:        now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"member",
	).SetMap(map[string]interface{}{
		"member_id":       m.ID,
		"user_id":         m.UserID,
		"team_id":         m.TeamID,
		"is_leader":       m.IsLeader,
		"invite_sent":     m.InviteSent,
		"invite_accepted": m.InviteAccepted,
		"created":         now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return m, errors.Wrapf(err, "inserting team member: %v", err)
	}

	return m, nil
}

func RetrieveMembers(ctx context.Context, repo *database.Repository, tid string) ([]Member, error) {
	var m []Member

	if _, err := uuid.Parse(tid); err != nil {
		return nil, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"member_id",
		"user_id",
		"team_id",
		"is_leader",
		"invite_sent",
		"invite_accepted",
		"created",
	).From(
		"team",
	).Where(sq.Eq{"team_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &m, q, tid); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return m, nil
}

func RetrieveMember(ctx context.Context, repo *database.Repository, tid, uid string) (Member, error) {
	var m Member

	if _, err := uuid.Parse(tid); err != nil {
		return m, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"member_id",
		"user_id",
		"team_id",
		"is_leader",
		"invite_sent",
		"invite_accepted",
		"created",
	).From(
		"team",
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

func Update(ctx context.Context, repo *database.Repository, tid string, update UpdateMember, uid string) error {
	p, err := RetrieveMember(ctx, repo, tid, uid)
	if err != nil {
		return err
	}

	if update.IsLeader != nil {
		p.IsLeader = *update.IsLeader
	}
	if update.InviteAccepted != nil {
		p.InviteAccepted = *update.InviteAccepted
	}

	stmt := repo.SQ.Update(
		"project",
	).SetMap(map[string]interface{}{
		"is_leader":       p.IsLeader,
		"invite_accepted": p.InviteAccepted,
	}).Where(sq.Eq{"team_id": tid, "user_id": uid})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating project")
	}

	return nil
}

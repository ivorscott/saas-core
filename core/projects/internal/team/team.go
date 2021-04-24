package team

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"time"
)


var (
	ErrNotFound  = errors.New("team not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

// Create adds a new Team
func Create(ctx context.Context, repo *database.Repository, nt NewTeam, pid, uid string, now time.Time) (*Team, error) {
	t := Team{
		ID:       uuid.New().String(),
		Name:     nt.Name,
		LeaderID: uid,
		Projects: []string{pid},
		Created:  now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"team",
	).SetMap(map[string]interface{}{
		"team_id":   t.ID,
		"name":      t.Name,
		"leader_id": t.LeaderID,
		"projects":  pq.Array(t.Projects),
		"created":   now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "inserting team: %v", nt)
	}

	return &t, nil
}

func Retrieve(ctx context.Context, repo *database.Repository, pid string) (*Team, error) {
	var t Team

	if _, err := uuid.Parse(pid); err != nil {
		return nil, ErrInvalidID
	}

	project := []string{pid}

	stmt := repo.SQ.Select(
		"team_id",
		"leader_id",
		"name",
		"projects",
		"created",
	).From(
		"team",
	).Where("? <@ projects")

	q, args, err := stmt.ToSql()
	log.Println(ctx, q)
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.QueryRowContext(ctx, q,  pq.Array(project)).Scan(&t.ID, &t.LeaderID, &t.Name, (*pq.StringArray)(&t.Projects), &t.Created);err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

//func CreateTeamMember(ctx context.Context, repo *database.Repository, m NewMember, tid string, time time.Time) (*Member, error) {
//
//
//	log.Println(tid, ni)
//		return nil, nil
//}

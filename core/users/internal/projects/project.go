package projects

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/devpies/devpie-client-core/users/internal/platform/database"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	ErrNotFound  = errors.New("project not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Retrieve(ctx context.Context, repo *database.Repository, pid string) (ProjectCopy, error) {
	var p ProjectCopy

	if _, err := uuid.Parse(pid); err != nil {
		return p, ErrInvalidID
	}

	// TODO: Check ownership at Team and User level before granting access
	// TODO: Also allow site Admins (where role==admin in Auth0)

	stmt := repo.SQ.Select(
		"project_id",
		"name",
		"team_id",
		"user_id",
		"active",
		"public",
		"column_order",
		"updated_at",
		"created_at",
	).From(
		"projects",
	).Where(sq.Eq{"project_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return p, errors.Wrapf(err, "building query: %v", args)
	}

	row := repo.DB.QueryRowContext(ctx, q, pid)
	err = row.Scan(&p.ID, &p.Name, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}

	return p, nil
}

func Create(ctx context.Context, repo *database.Repository, p ProjectCopy) error {
	stmt := repo.SQ.Insert(
		"projects",
	).SetMap(map[string]interface{}{
		"project_id":   p.ID,
		"name":         p.Name,
		"team_id":      p.TeamID,
		"user_id":      p.UserID,
		"active":       p.Active,
		"public":       p.Public,
		"column_order": pq.Array(p.ColumnOrder),
		"updated_at":   p.UpdatedAt,
		"created_at":   p.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "inserting project: %v", p)
	}

	return nil
}

func Update(ctx context.Context, repo *database.Repository, pid string, update UpdateProjectCopy)  error {
	p, err := Retrieve(ctx, repo, pid)
	if err != nil {
		return err
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Active != nil {
		p.Active = *update.Active
	}
	if update.Public != nil {
		p.Public = *update.Public
	}
	if update.TeamID != nil {
		p.TeamID = *update.TeamID
	}
	if update.ColumnOrder != nil {
		p.ColumnOrder = update.ColumnOrder
	}

	stmt := repo.SQ.Update(
		"projects",
	).SetMap(map[string]interface{}{
		"name":         p.Name,
		"active":       p.Active,
		"public":       p.Public,
		"column_order": pq.Array(p.ColumnOrder),
		"team_id":      p.TeamID,
		"updated_at":   update.UpdatedAt,
	}).Where(sq.Eq{"project_id": pid})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating project")
	}

	return nil
}

func Delete(ctx context.Context, repo *database.Repository, pid string) error {
	if _, err := uuid.Parse(pid); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"projects",
	).Where(sq.Eq{"project_id": pid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting project %s", pid)
	}

	return nil
}

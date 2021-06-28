package projects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Error codes returned by failures to handle projects.
var (
	ErrNotFound  = errors.New("project not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

// ProjectQuerier describes the behavior required for executing Project related queries
type ProjectQuerier interface {
	Create(ctx context.Context, repo *database.Repository, p ProjectCopy) error
	Retrieve(ctx context.Context, repo database.Storer, pid string) (ProjectCopy, error)
	Update(ctx context.Context, repo database.Storer, pid string, update UpdateProjectCopy) error
	Delete(ctx context.Context, repo database.Storer, pid string) error
}

// Queries defines method implementations for interacting with the projects table
type Queries struct{}

// Retrieve retrieves a single project from the database
func (q *Queries) Retrieve(ctx context.Context, repo database.Storer, pid string) (ProjectCopy, error) {
	var p ProjectCopy

	if _, err := uuid.Parse(pid); err != nil {
		return p, ErrInvalidID
	}

	stmt := repo.Select(
		"project_id",
		"name",
		"prefix",
		"description",
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

	query, args, err := stmt.ToSql()
	if err != nil {
		return p, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	row := repo.QueryRowxContext(ctx, query, pid)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}

	return p, nil
}

// Create inserts a new project into the database
func (q *Queries) Create(ctx context.Context, repo *database.Repository, p ProjectCopy) error {
	stmt := repo.Insert(
		"projects",
	).SetMap(map[string]interface{}{
		"project_id":   p.ID,
		"name":         p.Name,
		"prefix":       p.Prefix,
		"description":  p.Description,
		"team_id":      p.TeamID,
		"user_id":      p.UserID,
		"active":       p.Active,
		"public":       p.Public,
		"column_order": pq.Array(p.ColumnOrder),
		"updated_at":   p.UpdatedAt,
		"created_at":   p.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

// Update modifies a project in the database
func (q *Queries) Update(ctx context.Context, repo database.Storer, pid string, update UpdateProjectCopy) error {
	p, err := q.Retrieve(ctx, repo, pid)
	if err != nil {
		return err
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Description != nil {
		p.Description = *update.Description
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

	stmt := repo.Update(
		"projects",
	).SetMap(map[string]interface{}{
		"name":         p.Name,
		"description":  p.Description,
		"active":       p.Active,
		"public":       p.Public,
		"column_order": pq.Array(p.ColumnOrder),
		"team_id":      p.TeamID,
		"updated_at":   update.UpdatedAt,
	}).Where(sq.Eq{"project_id": pid})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a project from the database
func (q *Queries) Delete(ctx context.Context, repo database.Storer, pid string) error {
	if _, err := uuid.Parse(pid); err != nil {
		return ErrInvalidID
	}

	stmt := repo.Delete(
		"projects",
	).Where(sq.Eq{"project_id": pid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

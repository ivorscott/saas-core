package projects

import (
	"time"
)

// ProjectCopy represents a project copied from the projects service, stored for redundancy
type ProjectCopy struct {
	ID          string    `db:"project_id" json:"projectId"`
	Name        string    `db:"name" json:"name"`
	Prefix      string    `db:"prefix" json:"prefix"`
	Description string    `db:"description" json:"description"`
	TeamID      string    `db:"team_id" json:"teamId"`
	UserID      string    `db:"user_id" json:"userId"`
	Active      bool      `db:"active" json:"active"`
	Public      bool      `db:"public" json:"public"`
	ColumnOrder []string  `db:"column_order" json:"columnOrder"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

// UpdateProjectCopy represents an update to a project copy
type UpdateProjectCopy struct {
	Name        *string   `json:"name"`
	Active      *bool     `json:"active"`
	Public      *bool     `json:"public"`
	TeamID      *string   `json:"teamId"`
	ColumnOrder []string  `json:"columnOrder"`
	Description *string   `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

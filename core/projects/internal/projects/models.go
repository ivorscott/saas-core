package projects

import (
	"time"
)

type Project struct {
	ID          string    `db:"project_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	TeamID      string    `db:"team_id" json:"teamId"`
	UserID      string    `db:"user_id" json:"userId"`
	Active      bool      `db:"active" json:"active"`
	Public      bool      `db:"public" json:"public"`
	ColumnOrder []string  `db:"column_order" json:"columnOrder"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type NewProject struct {
	Name   string `json:"name" validate:"required"`
	TeamID string `json:"teamId"`
}

type UpdateProject struct {
	Name        *string   `json:"name"`
	Active      *bool     `json:"active"`
	Public      *bool     `json:"public"`
	TeamID      *string   `json:"teamId"`
	ColumnOrder []string `json:"columnOrder"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

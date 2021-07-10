package teams

import (
	"time"
)

// Team represents a group of members
type Team struct {
	ID        string    `db:"team_id" json:"id"`
	Name      string    `db:"name" json:"name"`
	UserID    string    `db:"user_id" json:"userId"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// NewTeam represents a new team request
type NewTeam struct {
	Name      string `json:"name" validate:"required"`
	ProjectID string `json:"projectId" validate:"required"`
}

// UpdateTeam represents an update to a team
type UpdateTeam struct {
	Name      *string   `json:"name"`
	UserID    *string   `json:"userId"`
	UpdatedAt time.Time `json:"updatedAt"`
}

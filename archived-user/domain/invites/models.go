package invites

import "time"

// Invite represents a team invitation
type Invite struct {
	ID         string    `db:"invite_id" json:"id"`
	UserID     string    `db:"user_id" json:"userId"`
	TeamID     string    `db:"team_id" json:"teamId"`
	Read       bool      `db:"read" json:"read"`
	Accepted   bool      `db:"accepted" json:"accepted"`
	Expiration time.Time `db:"expiration" json:"expiration"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

// InviteEnhanced represents an invitation enhanced with team details
type InviteEnhanced struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	TeamID     string    `json:"teamId"`
	TeamName   string    `json:"teamName"`
	Read       bool      `json:"read"`
	Accepted   bool      `json:"accepted"`
	Expiration time.Time `json:"expiration"`
	UpdatedAt  time.Time `json:"updatedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

// NewInvite represents a new team invite request
type NewInvite struct {
	UserID string `json:"userId" validate:"required"`
	TeamID string `json:"teamId" validate:"required"`
}

// NewList represents a list of email addresses to be invited
type NewList struct {
	Emails []string `json:"emailList" validate:"required"`
}

// UpdateInvite represents an update to an invite
type UpdateInvite struct {
	Accepted bool `json:"accepted" validate:"required"`
}

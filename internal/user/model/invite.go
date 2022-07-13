package model

import "time"

// Invite represents a team invitation
type Invite struct {
	ID         string    `db:"invite_id" json:"id"`
	TenantID   string    `db:"tenant_id" json:"tenantId"`
	UserID     string    `db:"user_id" json:"userId"`
	Read       bool      `db:"read" json:"read"`
	Accepted   bool      `db:"accepted" json:"accepted"`
	Expiration time.Time `db:"expiration" json:"expiration"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

// NewInvite represents a new tenant invite request
type NewInvite struct {
	UserID string `json:"userId" validate:"required"`
}

// UpdateInvite represents an update to an invite
type UpdateInvite struct {
	Accepted bool `json:"accepted" validate:"required"`
}

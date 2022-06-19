package model

import "time"

// MembershipCopy represents a Membership from the users service.
type MembershipCopy struct {
	ID        string    `db:"membership_id" json:"membershipId"`
	TenantID  string    `db:"tenant_id" json:"tenantId"`
	UserID    string    `db:"user_id" json:"userId"`
	TeamID    string    `db:"team_id" json:"teamId"`
	Role      string    `db:"role" json:"role"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// UpdateMembershipCopy represents a Membership update from the users service.
type UpdateMembershipCopy struct {
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updatedAt"`
}

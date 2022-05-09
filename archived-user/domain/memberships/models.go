package memberships

import "time"

// Membership represents a membership to a team
type Membership struct {
	ID        string    `db:"membership_id" json:"id"`
	UserID    string    `db:"user_id" json:"userId"`
	TeamID    string    `db:"team_id" json:"teamId"`
	Role      string    `db:"role" json:"role"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// MembershipEnhanced represents a membership enhanced with user details
type MembershipEnhanced struct {
	ID        string    `db:"membership_id" json:"id"`
	UserID    string    `db:"user_id" json:"userId"`
	FirstName *string   `db:"first_name" json:"firstName"`
	LastName  *string   `db:"last_name" json:"lastName"`
	Picture   *string   `db:"picture" json:"picture"`
	Email     string    `db:"email" json:"email"`
	TeamID    string    `db:"team_id" json:"teamId"`
	Role      string    `db:"role" json:"role"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// NewMembership represents a new membership request
type NewMembership struct {
	UserID string `json:"userId" validate:"required"`
	TeamID string `json:"teamId" validate:"required"`
	Role   string `json:"role" validate:"required"`
}

// UpdateMembership represents an update to a membership
type UpdateMembership struct {
	Role      *string   `json:"role"`
	UpdatedAt time.Time `json:"updatedAt"`
}

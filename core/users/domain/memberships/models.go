package memberships

import "time"

type Membership struct {
	ID        string    `db:"membership_id" json:"id"`
	UserID    string    `db:"user_id" json:"userId"`
	TeamID    string    `db:"team_id" json:"teamId"`
	Role      string    `db:"role" json:"role"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	Created   time.Time `db:"created_at" json:"createdAt"`
}

type NewMembership struct {
	UserID string `json:"userId"`
	TeamID string `json:"teamId"`
	Role   string `json:"role"`
}

type UpdateMembership struct {
	Role      *string   `json:"role"`
	UpdatedAt time.Time `json:"updatedAt"`
}

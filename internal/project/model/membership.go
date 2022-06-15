package model

import "time"

type MembershipCopy struct {
	ID        string    `db:"membership_id" json:"membershipId"`
	UserID    string    `db:"user_id" json:"userId"`
	TeamID    string    `db:"team_id" json:"teamId"`
	Role      string    `db:"role" json:"role"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type UpdateMembershipCopy struct {
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updatedAt"`
}

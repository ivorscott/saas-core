package membership

import "time"

type Membership struct {
	ID      string    `db:"membership_id" json:"id"`
	UserID  string    `db:"user_id" json:"user_id"`
	TeamID  string    `db:"team_id" json:"team_id"`
	Role    string    `db:"role" json:"role"`
	Created time.Time `db:"created" json:"created"`
}

type NewMembership struct {
	UserID string `db:"user_id" json:"user_id"`
	TeamID string `db:"team_id" json:"team_id"`
	Role   string `db:"role" json:"role"`
}

type UpdateMembership struct {
	Role *string `db:"role" json:"role"`
}

type Role int

const (
	Admin = iota
	Editor
	Commenter
	Viewer
)

func (r Role) String() string {
	return [...]string{"admin", "editor", "commenter", "viewer"}[r]
}

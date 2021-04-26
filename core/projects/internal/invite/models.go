package invite

import "time"

type Invite struct {
	ID         string    `db:"invite_id" json:"id"`
	UserID     string    `db:"user_id" json:"userId"`
	TeamID     string    `db:"team_id" json:"teamId"`
	Read       bool      `db:"read" json:"read"`
	Accepted   bool      `db:"accepted" json:"accepted"`
	Expiration time.Time `db:"expiration" json:"expiration"`
	Created    time.Time `db:"created" json:"created"`
}

type NewInvite struct {
	UserID string `json:"userId"`
	TeamID string `json:"teamId"`
}

type NewList struct {
	Emails []string `json:"emailList"`
}

type UpdateInvite struct {
	Read       *bool      `json:"read"`
	Accepted   *bool      `json:"accepted"`
	Expiration *time.Time `json:"expiration"`
}

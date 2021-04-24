package team

import (
	"time"
)

type Team struct {
	ID string `db:"team_id" json:"id"`
	LeaderID string `db:"leader_id" json:"leaderId"`
	Name string	`db:"name" json:"name"`
	Projects []string `db:"projects" json:"projects"`
	Created time.Time `db:"created" json:"created"`
}

type NewTeam struct {
	Name string `db:"name" json:"name"`
}

type Member struct {
	ID string `db:"member_id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`
	TeamID string `db:"team_id" json:"team_id"`
	IsLeader bool `db:"is_leader" json:"isLeader"`
	InviteSent bool `db:"invite_sent" json:"inviteSent"`
	InviteAccepted bool `db:"invite_accepted" json:"inviteAccepted"`
	Created time.Time `db:"created" json:"created"`
}

type NewMember struct {
	ID string `db:"member_id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`
	TeamID string `db:"team_id" json:"team_id"`
}
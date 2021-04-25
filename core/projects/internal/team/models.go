package team

import (
	"time"
)

type Team struct {
	ID       string    `db:"team_id" json:"id"`
	LeaderID string    `db:"leader_id" json:"leaderId"`
	Name     string    `db:"name" json:"name"`
	Created  time.Time `db:"created" json:"created"`
}

type NewTeam struct {
	Name string `db:"name" json:"name"`
}

type UpdateTeam struct {
	Name *string `db:"name" json:"name"`
}

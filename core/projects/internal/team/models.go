package team

import (
	"time"
)

type Team struct {
	ID      string    `db:"team_id" json:"id"`
	Name    string    `db:"name" json:"name"`
	UserId  string    `db:"user_id" json:"userId"`
	Created time.Time `db:"created" json:"created"`
}

type NewTeam struct {
	Name string `json:"name" validate:"required"`
}

type UpdateTeam struct {
	Name   *string `json:"name"`
	UserId *string `json:"userId"`
}

package teams

import (
	"time"
)

type Team struct {
	ID        string    `db:"team_id" json:"id"`
	Name      string    `db:"name" json:"name"`
	UserId    string    `db:"user_id" json:"userId"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreateAt  time.Time `db:"created_at" json:"createdAt"`
}

type NewTeam struct {
	Name string `json:"name" validate:"required"`
}

type UpdateTeam struct {
	Name      *string   `json:"name"`
	UserId    *string   `json:"userId"`
	UpdatedAt time.Time `json:"updatedAt"`
}

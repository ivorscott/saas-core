package user

import (
	"time"
)

type User struct {
	ID            string    `db:"user_id" json:"id" `
	Auth0ID       string    `db:"auth0_id" json:"auth0Id" `
	Email         string    `db:"email" json:"email"`
	Created       time.Time `db:"created" json:"created"`
}

type NewUser struct {
	Auth0ID       string  `json:"auth0Id" `
	Email         string  `json:"email"`
}

type UpdateUser struct {
	Email *string `json:"email"`
}
package main

// ========================================
// User Models
// ========================================

import (
	"time"
)

// ========================================
// User represents a registered auth0 user
// ========================================

type User struct {
	ID            string    `db:"user_id" json:"id" `
	Auth0ID       string    `db:"auth0_id" json:"auth0Id" `
	Email         string    `db:"email" json:"email"`
	EmailVerified bool      `db:"email_verified" json:"emailVerified"`
	FirstName     *string   `db:"first_name" json:"firstName"`
	LastName      *string   `db:"last_name" json:"lastName"`
	Picture       *string   `db:"picture" json:"picture"`
	Locale        *string   `db:"locale" json:"locale"`
	Created       time.Time `db:"created" json:"created"`
}

type NewUser struct {
	ID 			  string  `json:"user_id"`
	Auth0ID       string  `json:"auth0Id" `
	Email         string  `json:"email"`
	EmailVerified bool    `json:"emailVerified"`
	FirstName     *string `json:"firstName"`
	LastName      *string `json:"lastName"`
	Picture       *string `json:"picture"`
	Locale        *string `json:"locale"`
}

type UpdateUser struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Picture   *string `json:"picture"`
	Locale    *string `json:"locale"`
}

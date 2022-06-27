package model

import "time"

// User represent a user profile
type User struct {
	ID            string    `db:"user_id" json:"id"`
	TenantID      string    `db:"tenant_id" json:"tenantId"`
	Email         string    `db:"email" json:"email"`
	EmailVerified bool      `db:"email_verified" json:"emailVerified"`
	FirstName     string    `db:"first_name" json:"firstName"`
	LastName      string    `db:"last_name" json:"lastName"`
	Picture       *string   `db:"picture" json:"picture"`
	Locale        *string   `db:"locale" json:"locale"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}

// NewUser represents a new user request
type NewUser struct {
	Company       string  `json:"company" validate:"required"`
	Email         string  `json:"email" validate:"required"`
	FirstName     string  `json:"firstName" validate:"required"`
	LastName      string  `json:"lastName"`
	EmailVerified bool    `json:"emailVerified"`
	Picture       *string `json:"picture"`
	Locale        *string `json:"locale"`
}

// NewAdminUser represents a new user request
type NewAdminUser struct {
	UserID        string    `json:"userId" validate:"required"`
	TenantID      string    `json:"tenantId" validate:"required"`
	Company       string    `json:"company" validate:"required"`
	Email         string    `json:"email" validate:"required"`
	FirstName     string    `json:"firstName" validate:"required"`
	LastName      string    `json:"lastName" validate:"required"`
	EmailVerified bool      `json:"emailVerified" validate:"required"`
	CreatedAt     time.Time `json:"createdAt"`
}

// UpdateUser represents an update to a user
type UpdateUser struct {
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Picture   *string   `json:"picture"`
	Locale    *string   `json:"locale"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Seats struct {
	MaxSeats  int `db:"max_seats" json:"maxSeats"`
	SeatsUsed int `db:"seats_used" json:"seatsUsed"`
}

type SeatsAvailableResult struct {
	MaxSeats       int `json:"maxSeats"`
	SeatsAvailable int `json:"seatsAvailable"`
}

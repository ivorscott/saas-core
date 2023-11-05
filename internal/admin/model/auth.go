// Package model provides data transfer objects and validation.
package model

import (
	"github.com/go-playground/validator/v10"
)

var authValidator *validator.Validate

func init() {
	v := NewValidator()
	authValidator = v
}

// AuthCredentials represent login credentials.
type AuthCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Validate validates the AuthCredentials.
func (a *AuthCredentials) Validate() error {
	return authValidator.Struct(a)
}

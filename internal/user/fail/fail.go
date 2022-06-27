// Package fail contains common known errors.
package fail

import "errors"

var (
	// ErrNotFound represents a resource not found.
	ErrNotFound = errors.New("not found")
	// ErrInvalidID represents an invalid UUID.
	ErrInvalidID = errors.New("id provided was not a valid UUID")
	// ErrInvalidEmail represents an invalid email.
	ErrInvalidEmail = errors.New("email was not valid")
	// ErrConnectionFailed represents a failed connection attempt.
	ErrConnectionFailed = errors.New("connection failed")
	// ErrUserAlreadyAdded represents a failed attempt to add the same user a second time.
	ErrUserAlreadyAdded = errors.New("user already added")
)

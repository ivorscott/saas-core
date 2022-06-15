// Package fail contains common known errors.
package fail

import "errors"

var (
	// ErrNotFound represents a resource not found.
	ErrNotFound = errors.New("not found")
	// ErrInvalidID represents an invalid UUID.
	ErrInvalidID = errors.New("id provided was not a valid UUID")
	// ErrNotAuthorized represents an unauthorized request error.
	ErrNotAuthorized = errors.New("not authorized")
	// ErrConnectionFailed represents a failed connection attempt.
	ErrConnectionFailed = errors.New("connection failed")
)

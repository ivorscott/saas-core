// Package fail contains common known errors.
package fail

import "errors"

var (
	// ErrNotFound represents a resource not found.
	ErrNotFound = errors.New("not found")
	// ErrInvalidID represents an invalid UUID.
	ErrInvalidID = errors.New("id provided was not a valid UUID")
	// ErrNoTenant represents a failure to retrieve the tenant.
	ErrNoTenant = errors.New("missing tenant id")
)

package tenant

import "errors"

var (
	// ErrTenantNotFound represents a tenant not found error.
	ErrTenantNotFound = errors.New("tenant not found")
	// ErrUnauthorized represents an unauthorized action error.
	ErrUnauthorized = errors.New("unauthorized action")
	// ErrUnexpectedFailure represents an unexpected error.
	ErrUnexpectedFailure = errors.New("unexpected failure")
)

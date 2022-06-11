package tenant

import "errors"

var (
	ErrTenantNotFound    = errors.New("tenant not found")
	ErrUnauthorized      = errors.New("unauthorized action")
	ErrUnexpectedFailure = errors.New("unexpected failure")
)

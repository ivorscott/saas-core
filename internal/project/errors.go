package project

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidID        = errors.New("id provided was not a valid UUID")
	ErrNotAuthorized    = errors.New("user does not have correct membership")
	ErrConnectionFailed = errors.New("connection failed")
)

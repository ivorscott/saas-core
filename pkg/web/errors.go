package web

import (
	"errors"
)

// Error adds web information to request error.
type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

// Error returns the string error.
func (e *Error) Error() string {
	return e.Err.Error()
}

// NewRequestError is used when a known error condition is encountered.
func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

// shutdown is a type used to help with the graceful termination of the service.
type shutdown struct {
	Message string
}

func (s *shutdown) Error() string {
	return s.Message
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

// CtxErr returns an error for cases when values cannot be accessed from context.
func CtxErr() error {
	return NewShutdownError("cannot access values from context")
}

// IsShutdown checks to see if the shutdown error exists.
func IsShutdown(err error) bool {
	var targetErr *shutdown
	return errors.Is(err, targetErr)
}

// FieldError represents a request field error.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse represents the API error response.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

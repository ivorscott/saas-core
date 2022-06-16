package handler

import (
	"go.uber.org/zap"
)

type userService interface{}

// UserHandler handles the user requests.
type UserHandler struct {
	logger      *zap.Logger
	userService userService
}

// NewUserHandler returns a new user handler.
func NewUserHandler(
	logger *zap.Logger,
	userService userService,
) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

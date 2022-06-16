package service

import (
	"go.uber.org/zap"
)

type userRepository interface{}

// UserService manages the user business operations.
type UserService struct {
	logger   *zap.Logger
	userRepo userRepository
}

// NewUserHandler returns a new user handler.
func NewUserHandler(
	logger *zap.Logger,
	userRepo userRepository,
) *UserService {
	return &UserService{
		logger:   logger,
		userRepo: userRepo,
	}
}

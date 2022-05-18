package service

import (
	"context"
	"go.uber.org/zap"
)

// AuthService is responsible for managing authentication with Cognito.
type AuthService struct {
	logger *zap.Logger
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(logger *zap.Logger) *AuthService {
	return &AuthService{
		logger: logger,
	}
}

func (as *AuthService) Authenticate(ctx context.Context, email, password string) error {
	return nil
}

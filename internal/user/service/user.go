package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

type userRepository interface {
	Create(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error)
	RetrieveByEmail(ctx context.Context, email string) (model.User, error)
	RetrieveMe(ctx context.Context) (model.User, error)
}

// UserService manages the user business operations.
type UserService struct {
	logger   *zap.Logger
	userRepo userRepository
}

// NewUserService returns a new user service.
func NewUserService(
	logger *zap.Logger,
	userRepo userRepository,
) *UserService {
	return &UserService{
		logger:   logger,
		userRepo: userRepo,
	}
}

// AddSeat publishes a message to create a user in the identity service.
func (us *UserService) AddSeat(ctx context.Context, nu model.NewUser, now time.Time) error {
	return nil
}

func (us *UserService) RetrieveByEmail(ctx context.Context, email string) (model.User, error) {
	return us.userRepo.RetrieveByEmail(ctx, email)
}

func (us *UserService) RetrieveMe(ctx context.Context) (model.User, error) {
	return us.userRepo.RetrieveMe(ctx)
}

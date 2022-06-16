package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

type userRepository interface {
	Create(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error)
	RetrieveByEmail(email string) (model.User, error)
	RetrieveMe(ctx context.Context, uid string) (model.User, error)
	RetrieveMeByAuthID(ctx context.Context, aid string) (model.User, error)
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

func (us *UserService) Create(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error) {
	return us.userRepo.Create(ctx, nu, now)
}

func (us *UserService) RetrieveByEmail(email string) (model.User, error) {
	return us.userRepo.RetrieveByEmail(email)
}

func (us *UserService) RetrieveMe(ctx context.Context, uid string) (model.User, error) {
	return us.userRepo.RetrieveMe(ctx, uid)
}

func (us *UserService) RetrieveMeByAuthID(ctx context.Context, aid string) (model.User, error) {
	return us.userRepo.RetrieveMeByAuthID(ctx, aid)
}

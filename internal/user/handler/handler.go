package handler

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/user/model"
)

type userService interface {
	AddUser(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error)
	SeatsAvailable(ctx context.Context) (model.SeatsAvailableResult, error)
	List(ctx context.Context) ([]model.User, error)
	RetrieveByEmail(ctx context.Context, email string) (model.User, error)
	RetrieveMe(ctx context.Context) (model.User, error)
	RemoveUser(ctx context.Context, uid string) error
}

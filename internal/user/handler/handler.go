package handler

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/user/model"
)

type userService interface {
	AddSeat(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error)
	RetrieveByEmail(ctx context.Context, email string) (model.User, error)
	RetrieveMe(ctx context.Context, uid string) (model.User, error)
}

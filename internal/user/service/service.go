package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"time"
)

type inviteRepository interface {
	Create(ctx context.Context, ni model.NewInvite, now time.Time) (model.Invite, error)
	RetrieveInvite(ctx context.Context, iid string) (model.Invite, error)
	RetrieveInvites(ctx context.Context) ([]model.Invite, error)
	Update(ctx context.Context, update model.UpdateInvite, iid string, now time.Time) (model.Invite, error)
}

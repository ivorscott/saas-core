package service

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/user/model"

	"go.uber.org/zap"
)

// InviteService manages the user invite operations.
type InviteService struct {
	logger     *zap.Logger
	inviteRepo inviteRepository
}

// NewInviteService returns a new InviteService.
func NewInviteService(
	logger *zap.Logger,
	inviteRepo inviteRepository,
) *InviteService {
	return &InviteService{
		logger:     logger,
		inviteRepo: inviteRepo,
	}
}

// Create creates a new user invite.
func (is *InviteService) Create(ctx context.Context, ni model.NewInvite, now time.Time) (model.Invite, error) {
	return is.inviteRepo.Create(ctx, ni, now)
}

// RetrieveInvite retrieves a user invite.
func (is *InviteService) RetrieveInvite(ctx context.Context, iid string) (model.Invite, error) {
	return is.inviteRepo.RetrieveInvite(ctx, iid)
}

// RetrieveInvites retrieves user invites.
func (is *InviteService) RetrieveInvites(ctx context.Context) ([]model.Invite, error) {
	return is.inviteRepo.RetrieveInvites(ctx)
}

// Update updates a user invite.
func (is *InviteService) Update(ctx context.Context, update model.UpdateInvite, iid string, now time.Time) (model.Invite, error) {
	return is.inviteRepo.Update(ctx, update, iid, now)
}

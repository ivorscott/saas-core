package repository

import (
	"context"
	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

// InviteRepository manages invite data access.
type InviteRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewInviteRepository returns a new invite repository.
func NewInviteRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *InviteRepository {
	return &InviteRepository{
		logger: logger,
		pg:     pg,
	}
}

func (i InviteRepository) Create(ctx context.Context, ni model.NewInvite, now time.Time) (model.Invite, error) {
	//TODO implement me
	panic("implement me")
}

func (i InviteRepository) RetrieveInvite(ctx context.Context, uid string, iid string) (model.Invite, error) {
	//TODO implement me
	panic("implement me")
}

func (i InviteRepository) RetrieveInvites(ctx context.Context, uid string) ([]model.Invite, error) {
	//TODO implement me
	panic("implement me")
}

func (i InviteRepository) Update(ctx context.Context, update model.UpdateInvite, uid, iid string, now time.Time) (model.Invite, error) {
	//TODO implement me
	panic("implement me")
}

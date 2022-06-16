package repository

import (
	"context"
	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

// MembershipRepository manages membership data access.
type MembershipRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewMembershipRepository returns a new membership repository.
func NewMembershipRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *MembershipRepository {
	return &MembershipRepository{
		logger: logger,
		pg:     pg,
	}
}

func (mr *MembershipRepository) Create(ctx context.Context, nm model.NewMembership, now time.Time) (model.Membership, error) {
	//TODO implement me
	panic("implement me")
}

func (mr *MembershipRepository) RetrieveMemberships(ctx context.Context, uid, tid string) ([]model.MembershipEnhanced, error) {
	//TODO implement me
	panic("implement me")
}

func (mr *MembershipRepository) RetrieveMembership(ctx context.Context, uid, tid string) (model.Membership, error) {
	//TODO implement me
	panic("implement me")
}

func (mr *MembershipRepository) Update(ctx context.Context, tid string, update model.UpdateMembership, uid string, now time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (mr *MembershipRepository) Delete(ctx context.Context, tid, uid string) (string, error) {
	//TODO implement me
	panic("implement me")
}

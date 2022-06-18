package service

import (
	"context"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/pkg/msg"
	"go.uber.org/zap"
)

type membershipRepository interface {
	Create(ctx context.Context, nm model.MembershipCopy) error
	Update(ctx context.Context, mid string, update model.UpdateMembershipCopy) error
	Delete(ctx context.Context, mid string) error
	Retrieve(ctx context.Context, uid, tid string) (model.MembershipCopy, error)
	RetrieveByID(ctx context.Context, mid string) (model.MembershipCopy, error)
}

// MembershipService is responsible for managing redundant copy of membership data.
type MembershipService struct {
	logger *zap.Logger
	repo   membershipRepository
}

// NewMembershipService returns a new MembershipService.
func NewMembershipService(logger *zap.Logger, repo membershipRepository) *MembershipService {
	return &MembershipService{
		logger: logger,
		repo:   repo,
	}
}

func (ms *MembershipService) CreateMembershipFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalMembershipCreatedEvent(m)
	if err != nil {
		return err
	}

	nm := newMembership(event.Data)

	err = ms.repo.Create(ctx, nm)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MembershipService) UpdateMembershipFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalMembershipUpdatedEvent(m)
	if err != nil {
		return err
	}

	nm := newUpdateMembership(event.Data)

	err = ms.repo.Update(ctx, event.Data.MembershipID, nm)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MembershipService) DeleteMembershipFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalMembershipDeletedEvent(m)
	if err != nil {
		return err
	}

	err = ms.repo.Delete(ctx, event.Data.MembershipID)
	if err != nil {
		return err
	}
	return nil
}

func newUpdateMembership(data msg.MembershipUpdatedEventData) model.UpdateMembershipCopy {
	return model.UpdateMembershipCopy{
		Role:      data.Role,
		UpdatedAt: msg.ParseTime(data.UpdatedAt),
	}
}

func newMembership(data msg.MembershipCreatedEventData) model.MembershipCopy {
	return model.MembershipCopy{
		ID:        data.MembershipID,
		TenantID:  data.TenantID,
		CreatedAt: msg.ParseTime(data.CreatedAt),
		Role:      data.Role,
		TeamID:    data.TeamID,
		UpdatedAt: msg.ParseTime(data.UpdatedAt),
		UserID:    data.UserID,
	}
}

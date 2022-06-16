package handler

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type membershipService interface {
	Create(ctx context.Context, nm model.NewMembership, now time.Time) (model.Membership, error)
	RetrieveMemberships(ctx context.Context, uid, tid string) ([]model.MembershipEnhanced, error)
	RetrieveMembership(ctx context.Context, uid, tid string) (model.Membership, error)
	Update(ctx context.Context, tid string, update model.UpdateMembership, uid string, now time.Time) error
	Delete(ctx context.Context, tid, uid string) (string, error)
}

// MembershipHandler handles the membership requests.
type MembershipHandler struct {
	logger            *zap.Logger
	membershipService membershipService
}

// NewMembershipHandler returns a new membership handler.
func NewMembershipHandler(
	logger *zap.Logger,
	membershipService membershipService,
) *MembershipHandler {
	return &MembershipHandler{
		logger:            logger,
		membershipService: membershipService,
	}
}

func (mh *MembershipHandler) RetrieveMemberships(w http.ResponseWriter, r *http.Request) error {
	return nil
}

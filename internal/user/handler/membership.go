package handler

import (
	"context"
	"fmt"
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type membershipService interface {
	Create(ctx context.Context, nm model.NewMembership, now time.Time) (model.Membership, error)
	RetrieveMemberships(ctx context.Context, tid string) ([]model.MembershipEnhanced, error)
	RetrieveMembership(ctx context.Context, tid string) (model.Membership, error)
	Update(ctx context.Context, tid string, update model.UpdateMembership, now time.Time) error
	Delete(ctx context.Context, tid string) (string, error)
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

// RetrieveMemberships retrieves all memberships for the authenticated user.
func (mh *MembershipHandler) RetrieveMemberships(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	ms, err := mh.membershipService.RetrieveMemberships(r.Context(), tid)
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve memberships: %w", err)
		}
	}

	return web.Respond(r.Context(), w, ms, http.StatusOK)
}

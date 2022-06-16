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

// RetrieveMemberships retrieves all memberships for the authenticated user.
func (mh *MembershipHandler) RetrieveMemberships(w http.ResponseWriter, r *http.Request) error {
	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	tid := chi.URLParam(r, "tid")

	ms, err := mh.membershipService.RetrieveMemberships(r.Context(), values.Metadata.UserID, tid)
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
		}
		return fmt.Errorf("failed to retrieve memberships: %w", err)
	}

	return web.Respond(r.Context(), w, ms, http.StatusOK)
}

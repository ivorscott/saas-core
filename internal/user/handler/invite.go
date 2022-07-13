package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type inviteService interface {
	Create(ctx context.Context, ni model.NewInvite, now time.Time) (model.Invite, error)
	RetrieveInvite(ctx context.Context, iid string) (model.Invite, error)
	RetrieveInvites(ctx context.Context) ([]model.Invite, error)
	Update(ctx context.Context, update model.UpdateInvite, iid string, now time.Time) (model.Invite, error)
}

// InviteHandler handles the example requests.
type InviteHandler struct {
	logger        *zap.Logger
	inviteService inviteService
}

// NewInviteHandler returns a new invite handler.
func NewInviteHandler(
	logger *zap.Logger,
	inviteService inviteService,
) *InviteHandler {
	return &InviteHandler{
		logger:        logger,
		inviteService: inviteService,
	}
}

// CreateInvite sends new team invitations.
func (ih *InviteHandler) CreateInvite(w http.ResponseWriter, r *http.Request) error {
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// RetrieveInvites returns invitations for the authenticated user.
func (ih *InviteHandler) RetrieveInvites(w http.ResponseWriter, r *http.Request) error {
	var (
		is  []model.Invite
		err error
	)

	is, err = ih.inviteService.RetrieveInvites(r.Context())
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve invites: %w", err)
		}
	}

	return web.Respond(r.Context(), w, is, http.StatusOK)
}

// UpdateInvite updates an existing invitation.
func (ih *InviteHandler) UpdateInvite(w http.ResponseWriter, r *http.Request) error {
	var (
		update model.UpdateInvite
		err    error
	)

	if err = web.Decode(r, &update); err != nil {
		return err
	}

	iid := chi.URLParam(r, "iid")

	iv, err := ih.inviteService.Update(r.Context(), update, iid, time.Now())
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to update invite: %w", err)
		}
	}

	return web.Respond(r.Context(), w, iv, http.StatusOK)
}

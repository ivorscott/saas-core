package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

type teamRepository interface {
	Create(ctx context.Context, nt model.NewTeam, uid string, now time.Time) (model.Team, error)
	Retrieve(ctx context.Context, tid string) (model.Team, error)
	List(ctx context.Context, uid string) ([]model.Team, error)
}

type inviteRepository interface {
	Create(ctx context.Context, ni model.NewInvite, now time.Time) (model.Invite, error)
	RetrieveInvite(ctx context.Context, uid string, iid string) (model.Invite, error)
	RetrieveInvites(ctx context.Context, uid string) ([]model.Invite, error)
	Update(ctx context.Context, update model.UpdateInvite, uid, iid string, now time.Time) (model.Invite, error)
}

// TeamService manages the team business operations.
type TeamService struct {
	logger     *zap.Logger
	teamRepo   teamRepository
	inviteRepo inviteRepository
}

// NewTeamService returns a new team handler.
func NewTeamService(
	logger *zap.Logger,
	teamRepo teamRepository,
	inviteRepo inviteRepository,

) *TeamService {
	return &TeamService{
		logger:     logger,
		teamRepo:   teamRepo,
		inviteRepo: inviteRepo,
	}
}

func (ts *TeamService) Create(ctx context.Context, nt model.NewTeam, uid string, now time.Time) (model.Team, error) {
	return ts.teamRepo.Create(ctx, nt, uid, now)
}

func (ts *TeamService) Retrieve(ctx context.Context, tid string) (model.Team, error) {
	return ts.Retrieve(ctx, tid)
}

func (ts *TeamService) List(ctx context.Context, uid string) ([]model.Team, error) {
	return ts.teamRepo.List(ctx, uid)
}

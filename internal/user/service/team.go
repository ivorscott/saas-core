package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

type teamRepository interface {
	Create(ctx context.Context, nt model.NewTeam, now time.Time) (model.Team, error)
	Retrieve(ctx context.Context, tid string) (model.Team, error)
	List(ctx context.Context, uid string) ([]model.Team, error)
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

func (ts *TeamService) Create(ctx context.Context, nt model.NewTeam, now time.Time) (model.Team, error) {
	return ts.teamRepo.Create(ctx, nt, now)
}

func (ts *TeamService) Retrieve(ctx context.Context, tid string) (model.Team, error) {
	return ts.teamRepo.Retrieve(ctx, tid)
}

func (ts *TeamService) List(ctx context.Context, uid string) ([]model.Team, error) {
	return ts.teamRepo.List(ctx, uid)
}

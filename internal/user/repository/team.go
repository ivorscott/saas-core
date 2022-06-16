package repository

import (
	"context"
	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

// TeamRepository manages team data access.
type TeamRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewTeamRepository returns a new team repository.
func NewTeamRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *TeamRepository {
	return &TeamRepository{
		logger: logger,
		pg:     pg,
	}
}

func (tr *TeamRepository) Create(ctx context.Context, nt model.NewTeam, uid string, now time.Time) (model.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (tr *TeamRepository) Retrieve(ctx context.Context, tid string) (model.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (tr *TeamRepository) List(ctx context.Context, uid string) ([]model.Team, error) {
	//TODO implement me
	panic("implement me")
}

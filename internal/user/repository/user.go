package repository

import (
	"context"
	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

// UserRepository manages user data access.
type UserRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewUserRepository returns a new user repository.
func NewUserRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *UserRepository {
	return &UserRepository{
		logger: logger,
		pg:     pg,
	}
}

func (ur *UserRepository) Create(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (ur *UserRepository) RetrieveByEmail(email string) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (ur *UserRepository) RetrieveMe(ctx context.Context, uid string) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (ur *UserRepository) RetrieveMeByAuthID(ctx context.Context, aid string) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

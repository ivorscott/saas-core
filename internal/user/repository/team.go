package repository

import (
	"github.com/devpies/saas-core/internal/project/db"
	"go.uber.org/zap"
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

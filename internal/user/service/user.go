package service

import (
	"context"

	"github.com/devpies/saas-core/internal/user/model"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"go.uber.org/zap"
)

// UserService is responsible for managing users.
type UserService struct {
	logger        *zap.Logger
	userPoolID    string
	cognitoClient cognitoClient
}

type cognitoClient interface {
	AdminCreateUser(
		ctx context.Context,
		params *cognitoidentityprovider.AdminCreateUserInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

// NewUserService returns a new user service.
func NewUserService(logger *zap.Logger, userPoolID string, cognitoClient cognitoClient) *UserService {
	return &UserService{
		logger:        logger,
		userPoolID:    userPoolID,
		cognitoClient: cognitoClient,
	}
}

// CreateTenantUser creates a new user.
func (rs *UserService) CreateTenantUser(ctx context.Context, user model.NewUser) error {
	return nil
}

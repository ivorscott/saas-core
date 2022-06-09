package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

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

// CreateTenantUserFromMessage creates a new user from a NATS Message.
func (rs *UserService) CreateTenantUserFromMessage(ctx context.Context, message interface{}) error {
	data, err := msg.Bytes(message)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	m, err := msg.UnmarshalTenantRegisteredEvent(data)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	rs.logger.Info(fmt.Sprintf("%+v", m))
	return nil
}

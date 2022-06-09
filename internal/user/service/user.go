package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"net/http"

	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"go.uber.org/zap"
)

// UserService is responsible for managing users.
type UserService struct {
	logger        *zap.Logger
	cognitoClient cognitoClient
}

type cognitoClient interface {
	AdminCreateUser(
		ctx context.Context,
		params *cognitoidentityprovider.AdminCreateUserInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

// NewUserService returns a new user service.
func NewUserService(logger *zap.Logger, cognitoClient cognitoClient) *UserService {
	return &UserService{
		logger:        logger,
		cognitoClient: cognitoClient,
	}
}

// CreateTenantUserFromMessage creates a new user from a NATS Message.
func (rs *UserService) CreateTenantUserFromMessage(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	event, err := msg.UnmarshalTenantRegisteredEvent(m)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	d := event.Data

	_, err = rs.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(d.UserPoolID),
		Username:   aws.String(d.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:tenant-id"), Value: aws.String(d.ID)},
			{Name: aws.String("custom:company-name"), Value: aws.String(d.Company)},
			{Name: aws.String("email"), Value: aws.String(d.Email)},
			{Name: aws.String("email-verified"), Value: aws.String("true")},
		},
	})
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	rs.logger.Info("successfully added user")
	return nil
}

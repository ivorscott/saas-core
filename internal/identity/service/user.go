package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/devpies/saas-core/pkg/msg"

	"github.com/aws/aws-sdk-go-v2/aws"
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

// CreateTenantUserFromEvent creates a new user from a NATS Message.
func (rs *UserService) CreateTenantUserFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}

	event, err := msg.UnmarshalTenantRegisteredEvent(m)
	if err != nil {
		return err
	}

	d := event.Data

	_, err = rs.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(d.UserPoolID),
		Username:   aws.String(d.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:tenant-id"), Value: aws.String(d.ID)},
			{Name: aws.String("custom:account-owner"), Value: aws.String("1")},
			{Name: aws.String("custom:company-name"), Value: aws.String(formatPath(d.Company))},
			{Name: aws.String("custom:full-name"), Value: aws.String(fmt.Sprintf("%s %s", d.FirstName, d.LastName))},
			{Name: aws.String("email"), Value: aws.String(d.Email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	})
	if err != nil {
		rs.logger.Error("failed to add user", zap.Error(err))
		return err
	}
	rs.logger.Info("successfully added user")
	return nil
}

func formatPath(company string) string {
	return strings.ToLower(strings.Replace(company, " ", "", -1))
}

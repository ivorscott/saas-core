package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type publisher interface {
	Publish(subject string, message []byte)
}

// UserService is responsible for managing users.
type UserService struct {
	logger        *zap.Logger
	js            publisher
	cognitoClient cognitoClient
}

type cognitoClient interface {
	AdminCreateUser(
		ctx context.Context,
		params *cognitoidentityprovider.AdminCreateUserInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

// NewUserService returns a new user service.
func NewUserService(logger *zap.Logger, js publisher, cognitoClient cognitoClient) *UserService {
	return &UserService{
		logger:        logger,
		js:            js,
		cognitoClient: cognitoClient,
	}
}

// CreateTenantIdentityFromEvent creates a new identity from an event.
func (us *UserService) CreateTenantIdentityFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}

	event, err := msg.UnmarshalTenantRegisteredEvent(m)
	if err != nil {
		return err
	}

	d := event.Data

	user, err := us.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(d.UserPoolID),
		Username:   aws.String(d.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:tenant-id"), Value: aws.String(d.TenantID)},
			{Name: aws.String("custom:account-owner"), Value: aws.String("1")},
			{Name: aws.String("custom:company-name"), Value: aws.String(d.Company)},
			{Name: aws.String("custom:full-name"), Value: aws.String(fmt.Sprintf("%s %s", d.FirstName, d.LastName))},
			{Name: aws.String("email"), Value: aws.String(d.Email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	})
	if err != nil {
		us.logger.Error("failed to add user", zap.Error(err))
		return err
	}
	us.logger.Info("successfully added user")

	e, err := newIdentityCreatedEvent(ctx, user.User, d)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	us.js.Publish(msg.SubjectTenantIdentityCreated, bytes)
	return nil
}

func newIdentityCreatedEvent(
	ctx context.Context,
	user *types.UserType,
	data msg.TenantRegisteredEventData,
) (msg.TenantIdentityCreatedEvent, error) {
	var event msg.TenantIdentityCreatedEvent

	values, ok := web.FromContext(ctx)
	if !ok {
		return event, web.CtxErr()
	}

	var userID string
	for _, v := range user.Attributes {
		if v.Name != nil && *v.Name == "sub" {
			userID = *v.Value
			break
		}
	}

	event = msg.TenantIdentityCreatedEvent{
		Type: msg.TypeTenantIdentityCreated,
		Data: msg.TenantIdentityCreatedEventData{
			TenantID:  data.TenantID,
			UserID:    userID,
			Company:   data.Company,
			Email:     data.Email,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Plan:      data.Plan,
			CreatedAt: user.UserCreateDate.UTC().String(),
		},
		Metadata: msg.Metadata{
			TraceID: values.TraceID,
			UserID:  values.UserID,
		},
	}
	return event, nil
}

//go:generate mockery --quiet --all --dir . --case snake --output ../mocks --exported
package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/devpies/saas-core/internal/identity/mocks"
	"github.com/devpies/saas-core/internal/identity/service"
	"github.com/devpies/saas-core/pkg/msg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestUserService_CreateTenantUserFromMessage(t *testing.T) {
	var testCtx = context.Background()

	event := msg.TenantRegisteredEvent{
		Type: msg.TypeTenantRegistered,
		Metadata: msg.Metadata{
			TraceID: "123",
			UserID:  "123",
		},
		Data: msg.TenantRegisteredEventData{
			ID:         "tenant-id",
			Email:      "tenant-email",
			FirstName:  "first-name",
			LastName:   "last-name",
			Company:    "tenant-company",
			UserPoolID: "user-pool-id",
		},
	}
	eventBytes, _ := event.Marshal()
	message, _ := msg.UnmarshalMsg(eventBytes)

	fullName := fmt.Sprintf("%s %s", event.Data.FirstName, event.Data.LastName)

	input := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(event.Data.UserPoolID),
		Username:   aws.String(event.Data.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:tenant-id"), Value: aws.String(event.Data.ID)},
			{Name: aws.String("custom:account-owner"), Value: aws.String("1")},
			{Name: aws.String("custom:company-name"), Value: aws.String(event.Data.Company)},
			{Name: aws.String("custom:full-name"), Value: aws.String(fullName)},
			{Name: aws.String("email"), Value: aws.String(event.Data.Email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	}

	t.Run("400 error on create tenant", func(t *testing.T) {
		tests := []struct {
			name         string
			message      interface{}
			expectations func(deps userServiceDeps)
			err          string
		}{
			{
				name:         "failed to serialize message",
				message:      "invalid-message",
				expectations: func(deps userServiceDeps) {},
				err:          "not a message",
			},
			{
				name: fmt.Sprintf("failed to deserialize %s event", msg.TenantRegistered),
				message: &msg.Msg{
					Type:     msg.TenantRegistered,
					Metadata: msg.Metadata{TraceID: "123", UserID: "123"},
					Data:     "not event data",
				},
				expectations: func(deps userServiceDeps) {},
				err:          "json: cannot unmarshal string into Go struct field TenantRegisteredEvent.data of type msg.TenantRegisteredEventData",
			},
			{
				name:    "failed to create tenant user",
				message: &message,
				expectations: func(deps userServiceDeps) {
					deps.cognitoClient.On("AdminCreateUser", mock.AnythingOfType("*context.emptyCtx"), input).Return(nil, assert.AnError)
				},
				err: "assert.AnError general error for testing",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				svc, deps := setupUserServiceTest()

				tc.expectations(deps)

				err := svc.CreateTenantUserFromEvent(testCtx, tc.message)

				assert.Equal(t, tc.err, err.Error())
				deps.cognitoClient.AssertExpectations(t)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		svc, deps := setupUserServiceTest()
		output := &cognitoidentityprovider.AdminCreateUserOutput{}
		deps.cognitoClient.On("AdminCreateUser", mock.AnythingOfType("*context.emptyCtx"), input).Return(output, nil)

		err := svc.CreateTenantUserFromEvent(testCtx, &message)

		assert.Nil(t, err)
		deps.cognitoClient.AssertExpectations(t)
	})
}

type userServiceDeps struct {
	cognitoClient *mocks.CognitoClient
}

func setupUserServiceTest() (*service.UserService, userServiceDeps) {
	logger := zap.NewNop()
	cognitoClient := &mocks.CognitoClient{}
	return service.NewUserService(logger, cognitoClient), userServiceDeps{cognitoClient: cognitoClient}
}

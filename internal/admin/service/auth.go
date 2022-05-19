package service

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/devpies/core/internal/admin"
	"github.com/devpies/core/internal/admin/model"
	"go.uber.org/zap"
)

// AuthService is responsible for managing authentication with Cognito.
type AuthService struct {
	logger        *zap.Logger
	config        admin.Config
	cognitoClient cognitoClient
}

type cognitoClient interface {
	AdminInitiateAuth(ctx context.Context, params *cognitoidentityprovider.AdminInitiateAuthInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminInitiateAuthOutput, error)
}

var ErrMissingCognito = errors.New("missing cognito context")

// NewAuthService creates a new instance of AuthService.
func NewAuthService(logger *zap.Logger, config admin.Config, cognitoClient cognitoClient) *AuthService {
	if config.Cognito.AppClientID != "" && config.Cognito.UserPoolClientID != "" {
		return &AuthService{
			logger:        logger,
			config:        config,
			cognitoClient: cognitoClient,
		}
	}
	logger.Fatal("", zap.Error(ErrMissingCognito))
	return nil
}

func (as *AuthService) Authenticate(ctx context.Context, credentials model.AuthCredentials) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {
	signInInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:       "ADMIN_USER_PASSWORD_AUTH",
		ClientId:       aws.String(as.config.Cognito.AppClientID),
		UserPoolId:     aws.String(as.config.Cognito.UserPoolClientID),
		AuthParameters: map[string]string{"USERNAME": credentials.Email, "PASSWORD": credentials.Password},
	}

	return as.cognitoClient.AdminInitiateAuth(ctx, signInInput)
}

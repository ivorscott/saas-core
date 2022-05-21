package service

import (
	"context"
	"errors"

	"github.com/devpies/core/internal/adminapi/config"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

// AuthService is responsible for managing authentication with Cognito.
type AuthService struct {
	logger        *zap.Logger
	config        config.Config
	cognitoClient cognitoClient
}

type cognitoClient interface {
	AdminInitiateAuth(
		ctx context.Context,
		params *cognitoidentityprovider.AdminInitiateAuthInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminInitiateAuthOutput, error)
	AdminRespondToAuthChallenge(
		ctx context.Context,
		params *cognitoidentityprovider.AdminRespondToAuthChallengeInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminRespondToAuthChallengeOutput, error)
}

var ErrMissingCognito = errors.New("missing cognito context")

// NewAuthService creates a new instance of AuthService.
func NewAuthService(logger *zap.Logger, config config.Config, cognitoClient cognitoClient) *AuthService {
	return &AuthService{
		logger:        logger,
		config:        config,
		cognitoClient: cognitoClient,
	}
}

func (as *AuthService) Authenticate(ctx context.Context, email, password string) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {
	signInInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:       "ADMIN_USER_PASSWORD_AUTH",
		ClientId:       aws.String(as.config.Cognito.AppClientID),
		UserPoolId:     aws.String(as.config.Cognito.UserPoolClientID),
		AuthParameters: map[string]string{"USERNAME": email, "PASSWORD": password},
	}

	return as.cognitoClient.AdminInitiateAuth(ctx, signInInput)
}

func (as *AuthService) RespondToNewPasswordRequiredChallenge(ctx context.Context, email, password string, session string) (*cognitoidentityprovider.AdminRespondToAuthChallengeOutput, error) {
	params := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		ChallengeName: "NEW_PASSWORD_REQUIRED",
		ClientId:      aws.String(as.config.Cognito.AppClientID),
		UserPoolId:    aws.String(as.config.Cognito.UserPoolClientID),
		ChallengeResponses: map[string]string{
			"USERNAME":     email,
			"NEW_PASSWORD": password,
		},
		Session: aws.String(session),
	}
	return as.cognitoClient.AdminRespondToAuthChallenge(ctx, params)
}

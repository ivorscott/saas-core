package handler

import (
	"context"

	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type authService interface {
	Authenticate(ctx context.Context, email, password string) (*cip.AdminInitiateAuthOutput, error)
	CreateUserSession(ctx context.Context, token []byte) error
	CreatePasswordChallengeSession(ctx context.Context)
	RespondToNewPasswordRequiredChallenge(ctx context.Context, email, password string, session string) (*cip.AdminRespondToAuthChallengeOutput, error)
}

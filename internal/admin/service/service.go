package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type cognitoClient interface {
	AdminInitiateAuth(
		ctx context.Context,
		params *cognitoidentityprovider.AdminInitiateAuthInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminInitiateAuthOutput, error)
	AdminRespondToAuthChallenge(
		ctx context.Context,
		params *cognitoidentityprovider.AdminRespondToAuthChallengeInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminRespondToAuthChallengeOutput, error)
	AdminCreateUser(ctx context.Context, params *cognitoidentityprovider.AdminCreateUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

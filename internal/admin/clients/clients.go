package clients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func generateAccessToken(
	ctx context.Context,
	cognitoClient *cip.Client,
	credentials cognitoCredentials,
) (
	*string,
	error,
) {
	signInInput := &cip.AdminInitiateAuthInput{
		AuthFlow:   "ADMIN_USER_PASSWORD_AUTH",
		ClientId:   aws.String(credentials.cognitoClientID),
		UserPoolId: aws.String(credentials.userPoolID),
		AuthParameters: map[string]string{
			"USERNAME": credentials.m2mClientKey,
			"PASSWORD": credentials.m2mClientSecret,
		},
	}
	result, err := cognitoClient.AdminInitiateAuth(ctx, signInInput)
	if err != nil {
		return nil, err
	}
	return result.AuthenticationResult.AccessToken, nil
}

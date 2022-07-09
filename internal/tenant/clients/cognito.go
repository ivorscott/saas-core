package clients

import (
	"context"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// NewCognitoClient returns a new CognitoClient.
func NewCognitoClient(ctx context.Context, region string) *cip.Client {
	defaults, err := awsCfg.LoadDefaultConfig(ctx, awsCfg.WithRegion(region))
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return cip.NewFromConfig(defaults)
}

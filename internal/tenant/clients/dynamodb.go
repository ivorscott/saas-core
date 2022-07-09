package clients

import (
	"context"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewDynamoDBClient returns a client for DynamoDB.
func NewDynamoDBClient(ctx context.Context, region string) *dynamodb.Client {
	defaults, err := awsCfg.LoadDefaultConfig(ctx, awsCfg.WithRegion(region))
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return dynamodb.NewFromConfig(defaults)
}

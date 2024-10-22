// Package db creates a new database client.
package db

import (
	"context"

	"github.com/devpies/saas-core/internal/registration/config"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewDynamoDBClient returns a client for DynamoDB.
func NewDynamoDBClient(ctx context.Context, cfg config.Config) *dynamodb.Client {
	defaults, err := awsCfg.LoadDefaultConfig(ctx, awsCfg.WithRegion(cfg.Cognito.Region))
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return dynamodb.NewFromConfig(defaults)
}

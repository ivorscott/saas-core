package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devpies/saas-core/pkg/lambda/modifytoken/repository"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {
	var err error

	client := repository.NewDynamoRepository(ctx, event.Region)

	// Early return if tenant-id does not exist.
	// e.g., an M2M client application will not have a tenant id.
	if _, ok := event.Request.UserAttributes["custom:tenant-id"]; !ok {
		return event, nil
	}

	connections, err := client.FindTenantConnections(ctx, event.Request.UserAttributes["sub"])
	if err != nil {
		return event, err
	}

	bytes, err := json.Marshal(&connections)
	if err != nil {
		return event, err
	}

	event.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride = map[string]string{
		"custom:tenant-connections": string(bytes),
	}
	return event, nil
}

func main() {
	lambda.Start(handler)
}

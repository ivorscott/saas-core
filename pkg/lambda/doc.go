// Package lambda provides AWS Lambda functions required by the app.
//
// The modifytoken package modifies the id token before it is generated. It inserts a tenant connection mapping,
// containing all the tenants the user has access to.
// https://github.com/aws/aws-lambda-go/blob/main/events/README_Cognito_UserPools_PreTokenGen.md
// https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-pre-token-generation.html
//
// Build instructions:
// 1. cd pkg/lambda/modifytoken
// 1. GOARCH=amd64 GOOS=linux go build -o modifytoken
// 2. zip modifytoken.zip modifytoken
//
// The updatestatus package modifies the status of the tenant post sign up confirmation.
// https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-post-confirmation.html
// Build instructions:
// 1. cd pkg/lambda/updatestatus
// 1. GOARCH=amd64 GOOS=linux go build -o updatestatus
// 2. zip updatestatus.zip updatestatus
package lambda

package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/devpies/saas-core/internal/registration/config"
	"github.com/devpies/saas-core/internal/registration/model"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"go.uber.org/zap"
)

// IDPService manages the identity provider configuration.
type IDPService struct {
	logger        *zap.Logger
	config        config.Config
	cognitoClient cognitoClient
	authInfoRepo  authInfoRepository
	js            publisher
}

type cognitoClient interface {
	CreateUserPool(ctx context.Context, params *cognitoidentityprovider.CreateUserPoolInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateUserPoolOutput, error)
	CreateUserPoolClient(ctx context.Context, params *cognitoidentityprovider.CreateUserPoolClientInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateUserPoolClientOutput, error)
}

type authInfoRepository interface {
	InsertAuthInfo(ctx context.Context, info model.AuthInfo) error
	SelectAuthInfo(ctx context.Context, path string) (model.AuthInfo, error)
}

const (
	// UserPoolSiloed represents the user pool type for premium tenants.
	UserPoolSiloed = "siloed"
	// UserPoolPooled represents a shared user pool for basic tier tenants.
	UserPoolPooled = "pooled"
	// DefaultTenantPath represents the tenant path for basic tier tenants in a shared pool.
	DefaultTenantPath = "/app"
)

// NewIDPService returns a new IDPService.
func NewIDPService(logger *zap.Logger, config config.Config, cognitoClient cognitoClient, authInfoRepo authInfoRepository, js publisher) *IDPService {
	return &IDPService{
		logger:        logger,
		config:        config,
		cognitoClient: cognitoClient,
		authInfoRepo:  authInfoRepo,
		js:            js,
	}
}

// GetPlanBasedUserPool either retrieves an existing user pool (shared pool) or creates a new pool (siloed). In both cases,
// the user pool id is returned. If a new user pool is required a user pool client will also be created.
func (idps *IDPService) GetPlanBasedUserPool(ctx context.Context, tenant model.NewTenant, path string) (string, error) {
	var (
		poolType  = UserPoolPooled
		pathToUse = DefaultTenantPath
		err       error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return "", web.CtxErr()
	}

	if Plan(tenant.Plan) == PlanPremium {
		poolType = UserPoolSiloed
		pathToUse = path
	}

	// Fetch existing pool id and exit if one exists.
	existingPoolID, err := idps.fetchPoolID(ctx, pathToUse)
	if err != nil {
		return "", err
	}
	if existingPoolID != "" {
		return existingPoolID, nil
	}

	// Otherwise, create a user pool.
	poolName := fmt.Sprintf("siloed-%s", tenant.ID)
	userPool, err := idps.createUserPool(ctx, poolName, pathToUse)
	if err != nil {
		return "", err
	}
	// Create user pool client.
	userPoolClient, err := idps.createUserPoolClient(ctx, tenant.ID, userPool.Id)
	if err != nil {
		return "", err
	}
	// Store auth info for path.
	err = idps.storeForPath(ctx, pathToUse, poolType, userPool.Id, userPoolClient.UserPoolId)
	if err != nil {
		return "", err
	}
	// Publish siloed tenant config
	event := newCreateTenantSiloedEvent(values, tenant.Company, userPoolClient.ClientId, userPool.Id)
	bytes, err := event.Marshal()
	if err != nil {
		return "", nil
	}
	idps.js.Publish(msg.SubjectSiloed, bytes)
	return *userPool.Id, nil
}

func newCreateTenantSiloedEvent(values *web.Values, path string, clientID, userPoolID *string) msg.TenantSiloedEvent {
	return msg.TenantSiloedEvent{
		Metadata: msg.Metadata{
			TraceID: values.Metadata.TraceID,
			UserID:  values.Metadata.UserID,
		},
		Type: msg.TypeTenantSiloed,
		Data: msg.TenantSiloedEventData{
			TenantName:       strings.TrimPrefix(path, "/"),
			AppClientID:      *clientID,
			UserPoolID:       *userPoolID,
			DeploymentStatus: "provisioned",
		},
	}
}

func (idps *IDPService) fetchPoolID(ctx context.Context, path string) (string, error) {
	info, err := idps.authInfoRepo.SelectAuthInfo(ctx, path)
	if err != nil {
		return "", err
	}
	return info.UserPoolID, nil
}

func (idps *IDPService) storeForPath(ctx context.Context, path, poolType string, userPoolID, userPoolClientID *string) error {
	info := model.AuthInfo{
		TenantPath:       path,
		UserPoolID:       *userPoolID,
		UserPoolType:     poolType,
		UserPoolClientID: *userPoolClientID,
	}
	err := idps.authInfoRepo.InsertAuthInfo(ctx, info)
	if err != nil {
		return err
	}
	return nil
}

func (idps *IDPService) createUserPool(ctx context.Context, poolName, path string) (*types.UserPoolType, error) {
	protocol := "http"
	host := "localhost:3000"

	if idps.config.Web.Production {
		protocol = "https"
		host = "devpie.io"
	}

	url := fmt.Sprintf(`<a href="%s://%s/%s\">%s://%s/%s</a>`, protocol, host, path, protocol, host, path)

	// Email template for on-boarding tenant.
	emailTemplate := fmt.Sprintf(`<b>Welcome to the SaaS Application for EKS Workshop!</b> <br>
    <br>
    The URL for your application is here: %s. 
    <br>
    <br>
    Please note that it may take a few minutes to provision your tenant. If you get a 404 when hitting the link above
    please try again in a few minutes. You can also check the AWS CodePipeline project that's in your environment
    for status.
    <br>
    Your username is: <b>{username}</b>
    <br>
    Your temporary password is: <b>{####}</b>
    <br>`, url)

	input := cognitoidentityprovider.CreateUserPoolInput{
		PoolName: aws.String(poolName),
		AdminCreateUserConfig: &types.AdminCreateUserConfigType{
			AllowAdminCreateUserOnly: true,
			InviteMessageTemplate: &types.MessageTemplateType{
				EmailMessage: aws.String(emailTemplate),
				EmailSubject: aws.String("Temporary password for environment EKS SaaS Application"),
			},
		},
		UsernameAttributes: []types.UsernameAttributeType{"email"},
		Schema: []types.SchemaAttributeType{
			{
				AttributeDataType: "String",
				Name:              aws.String("email"),
				Required:          true,
				Mutable:           true,
			},
			{
				AttributeDataType: "String",
				Name:              aws.String("tenant-id"),
				Required:          false,
				Mutable:           false,
			},
			{
				AttributeDataType: "String",
				Name:              aws.String("company-name"),
				Required:          false,
				Mutable:           false,
			},
			{
				AttributeDataType: "String",
				Name:              aws.String("full-name"),
				Required:          false,
				Mutable:           true,
			},
		},
	}

	// Create a new Amazon Cognito user pool.
	output, err := idps.cognitoClient.CreateUserPool(ctx, &input)
	if err != nil {
		return nil, err
	}
	return output.UserPool, err
}

func (idps *IDPService) createUserPoolClient(ctx context.Context, tenantID string, userPoolID *string) (*types.UserPoolClientType, error) {
	input := cognitoidentityprovider.CreateUserPoolClientInput{
		ClientName: aws.String(tenantID),
		UserPoolId: userPoolID,
		ExplicitAuthFlows: []types.ExplicitAuthFlowsType{
			types.ExplicitAuthFlowsTypeAllowAdminUserPasswordAuth,
			types.ExplicitAuthFlowsTypeAllowUserSrpAuth,
			types.ExplicitAuthFlowsTypeAllowRefreshTokenAuth,
		},
		GenerateSecret:             false,
		PreventUserExistenceErrors: types.PreventUserExistenceErrorTypesEnabled,
		RefreshTokenValidity:       30,
		SupportedIdentityProviders: []string{"COGNITO"},
	}
	output, err := idps.cognitoClient.CreateUserPoolClient(ctx, &input)
	if err != nil {
		return nil, err
	}
	return output.UserPoolClient, nil
}

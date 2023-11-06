package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"go.uber.org/zap"
)

type registrationClient interface {
	Register(ctx context.Context, tenant []byte) (*http.Response, error)
}

// RegistrationService is responsible for triggering tenant registration.
type RegistrationService struct {
	logger           *zap.Logger
	sharedUserPoolID string
	cognitoClient    cognitoClient
	httpClient       registrationClient
}

// NewRegistrationService returns a new registration service.
func NewRegistrationService(logger *zap.Logger, sharedUserPoolID string, cognitoClient cognitoClient, httpClient registrationClient) *RegistrationService {
	return &RegistrationService{
		logger:           logger,
		sharedUserPoolID: sharedUserPoolID,
		cognitoClient:    cognitoClient,
		httpClient:       httpClient,
	}
}

// RegisterTenant sends new tenant to tenant registration microservice and returns a nil ErrorResponse on success.
func (rs *RegistrationService) RegisterTenant(ctx context.Context, newTenant model.NewTenant) (int, error) {
	var (
		resp *http.Response
		err  error
	)

	// Tenants default to a basic plan
	newTenant.Plan = "basic"

	data, err := json.Marshal(newTenant)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	resp, err = rs.httpClient.Register(ctx, data)
	if err != nil {
		rs.logger.Error("error performing request")
		return http.StatusInternalServerError, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		rs.logger.Error("error reading body")
		return http.StatusInternalServerError, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var webErrResp web.ErrorResponse
		err = json.Unmarshal(bodyBytes, &webErrResp)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = &web.Error{
			Err:    fmt.Errorf(webErrResp.Error),
			Status: resp.StatusCode,
			Fields: webErrResp.Fields,
		}
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}

// ResendTemporaryPassword resends the user a temporary password.
func (rs *RegistrationService) ResendTemporaryPassword(ctx context.Context, username string) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	return rs.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		Username:      aws.String(username),
		UserPoolId:    aws.String(rs.sharedUserPoolID),
		MessageAction: types.MessageActionTypeResend,
	})
}

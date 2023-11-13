package clients

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devpies/saas-core/pkg/web"

	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"go.uber.org/zap"
)

// HTTPSubscriptionClient manages calls to the registration service.
type HTTPSubscriptionClient struct {
	logger         *zap.Logger
	client         *http.Client
	cognito        *cip.Client
	credentials    cognitoCredentials
	serviceAddress string
	servicePort    string
}

type cognitoCredentials struct {
	cognitoClientID string
	userPoolID      string
	m2mClientKey    string
	m2mClientSecret string
}

// NewHTTPSubscriptionClient returns a new HttpSubscriptionClient.
func NewHTTPSubscriptionClient(
	logger *zap.Logger,
	serviceAddress string,
	servicePort string,
	cognitoClient *cip.Client,
	cognitoClientID string,
	userPoolID string,
	m2mClientKey string,
	m2mClientSecret string,
) *HTTPSubscriptionClient {
	return &HTTPSubscriptionClient{
		logger:         logger,
		client:         &http.Client{},
		cognito:        cognitoClient,
		serviceAddress: serviceAddress,
		servicePort:    servicePort,
		credentials: cognitoCredentials{
			cognitoClientID: cognitoClientID,
			userPoolID:      userPoolID,
			m2mClientKey:    m2mClientKey,
			m2mClientSecret: m2mClientSecret,
		},
	}
}

// GetSubscriptionInfo calls the subscription service over a http interface to retrieve a tenant's subscription information.
// The subscription service is a "public" service therefore a user token must be generated to authenticate.
func (h *HTTPSubscriptionClient) GetSubscriptionInfo(ctx context.Context, tenantID string) (*http.Response, error) {
	userToken, err := generateAccessToken(ctx, h.cognito, h.credentials)
	if err != nil {
		return nil, err
	}

	if userToken == nil {
		h.logger.Error("generated cognito user token is nil")
	}

	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}

	url := fmt.Sprintf("http://%s:%s/subscriptions/%s", h.serviceAddress, h.servicePort, tenantID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("TraceID", values.TraceID)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *userToken))

	client := &http.Client{}
	return client.Do(request)
}

// RefundUser makes a refund request to the subscription service.
func (h *HTTPSubscriptionClient) RefundUser(ctx context.Context, subID string) (*http.Response, error) {
	userToken, err := generateAccessToken(ctx, h.cognito, h.credentials)
	if err != nil {
		return nil, err
	}

	if userToken == nil {
		h.logger.Error("generated cognito user token is nil")
	}

	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}

	url := fmt.Sprintf("http://%s:%s/subscriptions/refund/%s", h.serviceAddress, h.servicePort, subID)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("TraceID", values.TraceID)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *userToken))

	client := &http.Client{}
	return client.Do(request)
}

// CancelSubscription makes a subscription cancellation request to the subscription service.
func (h *HTTPSubscriptionClient) CancelSubscription(ctx context.Context, subID string) (*http.Response, error) {
	userToken, err := generateAccessToken(ctx, h.cognito, h.credentials)
	if err != nil {
		return nil, err
	}

	if userToken == nil {
		h.logger.Error("generated cognito user token is nil")
	}

	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}

	url := fmt.Sprintf("http://%s:%s/subscriptions/cancel/%s", h.serviceAddress, h.servicePort, subID)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("TraceID", values.TraceID)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *userToken))

	client := &http.Client{}
	return client.Do(request)
}

package clients

import (
	"context"
	"fmt"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"net/http"

	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

// HTTPBillingClient manages calls to the registration service.
type HTTPBillingClient struct {
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

// NewHTTPBillingClient returns a new HttpBillingClient.
func NewHTTPBillingClient(
	logger *zap.Logger,
	serviceAddress string,
	servicePort string,
	cognitoClient *cip.Client,
	cognitoClientID string,
	userPoolID string,
	m2mClientKey string,
	m2mClientSecret string,
) *HTTPBillingClient {
	return &HTTPBillingClient{
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

// FindAllSubscriptions calls the billing service over a http interface to retrieve all subscriptions.
// The billing service is a "public" service therefore a user token must be generated to authenticate.
func (h *HTTPBillingClient) FindAllSubscriptions(ctx context.Context) (*http.Response, error) {
	userToken, err := generateAccessToken(ctx, h.cognito, h.credentials)
	if err != nil {
		return nil, err
	}

	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}

	url := fmt.Sprintf("http://%s:%s/subscriptions", h.serviceAddress, h.servicePort)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("TraceID", values.TraceID)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))

	client := &http.Client{}
	return client.Do(request)
}

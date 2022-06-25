package clients

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

// HTTPRegistrationClient manages calls to the registration service.
type HTTPRegistrationClient struct {
	client         *http.Client
	logger         *zap.Logger
	serviceAddress string
	servicePort    string
}

// NewHTTPRegistrationClient returns a new HttpRegistrationClient.
func NewHTTPRegistrationClient(logger *zap.Logger, serviceAddress string, servicePort string) *HTTPRegistrationClient {
	return &HTTPRegistrationClient{
		logger:         logger,
		client:         &http.Client{},
		serviceAddress: serviceAddress,
		servicePort:    servicePort,
	}
}

// Register calls the registration service over a http interface to register a new tenant.
func (h *HTTPRegistrationClient) Register(ctx context.Context, tenant []byte) (*http.Response, error) {
	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}

	payload := bytes.NewReader(tenant)
	url := fmt.Sprintf("http://%s:%s/registration/register", h.serviceAddress, h.servicePort)

	request, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("TraceID", values.TraceID)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", values.Token))

	client := &http.Client{}
	return client.Do(request)
}

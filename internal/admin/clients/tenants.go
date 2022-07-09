package clients

import (
	"context"
	"fmt"
	"github.com/devpies/saas-core/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

// HTTPTenantClient manages calls to the registration service.
type HTTPTenantClient struct {
	client         *http.Client
	logger         *zap.Logger
	serviceAddress string
	servicePort    string
}

// NewHTTPTenantClient returns a new HttpTenantClient.
func NewHTTPTenantClient(logger *zap.Logger, serviceAddress string, servicePort string) *HTTPTenantClient {
	return &HTTPTenantClient{
		logger:         logger,
		client:         &http.Client{},
		serviceAddress: serviceAddress,
		servicePort:    servicePort,
	}
}

// FindAllTenants calls the tenant service over a http interface to retrieve all tenants.
func (h *HTTPTenantClient) FindAllTenants(ctx context.Context) (*http.Response, error) {
	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}

	url := fmt.Sprintf("http://%s:%s/tenants", h.serviceAddress, h.servicePort)

	request, err := http.NewRequest(http.MethodGet, url, nil)
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

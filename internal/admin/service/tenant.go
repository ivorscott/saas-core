package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type tenantClient interface {
	FindAllTenants(ctx context.Context) (*http.Response, error)
}

type billingClient interface {
	FindAllSubscriptions(ctx context.Context) (*http.Response, error)
}

// TenantService manages the example business operations.
type TenantService struct {
	logger        *zap.Logger
	tenantClient  tenantClient
	billingClient billingClient
}

// NewTenantService returns a new example service.
func NewTenantService(
	logger *zap.Logger,
	tenantClient tenantClient,
	billingClient billingClient,
) *TenantService {
	return &TenantService{
		logger:        logger,
		tenantClient:  tenantClient,
		billingClient: billingClient,
	}
}

// ListTenants lists all tenants returned by the tenant microservice.
func (ts *TenantService) ListTenants(ctx context.Context) ([]model.Tenant, int, error) {
	var (
		resp    *http.Response
		tenants []model.Tenant
		err     error
	)

	resp, err = ts.tenantClient.FindAllTenants(ctx)
	if err != nil {
		ts.logger.Error("error performing request")
		return nil, http.StatusInternalServerError, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ts.logger.Error("error reading body")
		return nil, http.StatusInternalServerError, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var webErrResp web.ErrorResponse
		err = json.Unmarshal(bodyBytes, &webErrResp)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = &web.Error{
			Err:    fmt.Errorf(webErrResp.Error),
			Status: resp.StatusCode,
			Fields: webErrResp.Fields,
		}
		return nil, resp.StatusCode, err
	}

	err = json.Unmarshal(bodyBytes, &tenants)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return tenants, resp.StatusCode, nil
}

// ListTenantSubscriptions retrieves all the paid subscriptions from the billing service.
func (ts *TenantService) ListTenantSubscriptions(ctx context.Context) ([]model.Subscription, int, error) {
	var (
		resp          *http.Response
		subscriptions []model.Subscription
		err           error
	)

	resp, err = ts.billingClient.FindAllSubscriptions(ctx)
	if err != nil {
		ts.logger.Error("error performing request")
		return nil, http.StatusInternalServerError, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ts.logger.Error("error reading body")
		return nil, http.StatusInternalServerError, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var webErrResp web.ErrorResponse
		err = json.Unmarshal(bodyBytes, &webErrResp)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = &web.Error{
			Err:    fmt.Errorf(webErrResp.Error),
			Status: resp.StatusCode,
			Fields: webErrResp.Fields,
		}
		return nil, resp.StatusCode, err
	}

	err = json.Unmarshal(bodyBytes, &subscriptions)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return subscriptions, resp.StatusCode, nil
}

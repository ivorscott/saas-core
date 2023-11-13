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
	GetSubscriptionInfo(ctx context.Context, tenantID string) (*http.Response, error)
	RefundUser(ctx context.Context, subID string) (*http.Response, error)
	CancelSubscription(ctx context.Context, subID string) (*http.Response, error)
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

// GetSubscriptionInfo finds a tenant's subscription information.
func (ts *TenantService) GetSubscriptionInfo(ctx context.Context, tenantID string) (model.SubscriptionInfo, int, error) {
	var (
		resp         *http.Response
		subscription model.SubscriptionInfo
		err          error
	)

	resp, err = ts.billingClient.GetSubscriptionInfo(ctx, tenantID)
	if err != nil {
		ts.logger.Error("error performing request")
		return subscription, http.StatusInternalServerError, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ts.logger.Error("error reading body")
		return subscription, http.StatusInternalServerError, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var webErrResp web.ErrorResponse
		err = json.Unmarshal(bodyBytes, &webErrResp)
		if err != nil {
			return subscription, http.StatusInternalServerError, err
		}
		err = &web.Error{
			Err:    fmt.Errorf(webErrResp.Error),
			Status: resp.StatusCode,
			Fields: webErrResp.Fields,
		}
		return subscription, resp.StatusCode, err
	}

	err = json.Unmarshal(bodyBytes, &subscription)
	if err != nil {
		return subscription, http.StatusInternalServerError, err
	}

	return subscription, resp.StatusCode, nil
}

// RefundUser refunds the stripe user.
func (ts *TenantService) RefundUser(ctx context.Context, subID string) (int, error) {
	var (
		resp         *http.Response
		subscription model.SubscriptionInfo
		err          error
	)

	resp, err = ts.billingClient.RefundUser(ctx, subID)
	if err != nil {
		ts.logger.Error("error performing request")
		return http.StatusInternalServerError, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ts.logger.Error("error reading body")
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

	err = json.Unmarshal(bodyBytes, &subscription)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}

// CancelSubscription cancels the subscription for the stripe user.
func (ts *TenantService) CancelSubscription(ctx context.Context, tenantID string) (int, error) {
	var (
		resp         *http.Response
		subscription model.SubscriptionInfo
		err          error
	)

	resp, err = ts.billingClient.CancelSubscription(ctx, tenantID)
	if err != nil {
		ts.logger.Error("error performing request")
		return http.StatusInternalServerError, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ts.logger.Error("error reading body")
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

	err = json.Unmarshal(bodyBytes, &subscription)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}

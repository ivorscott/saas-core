package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/internal/admin/render"
)

type renderer interface {
	Template(w http.ResponseWriter, r *http.Request, page string, td *render.TemplateData, partials ...string) error
}

type tenantService interface {
	ListTenants(ctx context.Context) ([]model.Tenant, int, error)
	GetSubscriptionInfo(ctx context.Context, tenantID string) (model.SubscriptionInfo, int, error)
	CancelSubscription(ctx context.Context, subID string) (int, error)
	RefundUser(ctx context.Context, subID string) (int, error)
}

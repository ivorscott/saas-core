package service

import (
	"context"
	"fmt"

	"github.com/devpies/saas-core/internal/registration/model"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type publisher interface {
	Publish(subject string, message []byte)
}

// RegistrationService is responsible for managing tenant registration.
type RegistrationService struct {
	logger            *zap.Logger
	js                publisher
	tenantsStream     string
	defaultUserPoolID string
}

// Plan represents the type of subscription plan.
type Plan string

const (
	// PlanBasic represents the cheapest subscription plan offering.
	PlanBasic Plan = "basic"
	// PlanPremium represents the premium plan offering.
	PlanPremium Plan = "premium"
)

// NewRegistrationService returns a new registration service.
func NewRegistrationService(logger *zap.Logger, js publisher, defaultUserPoolID, tenantsStream string) *RegistrationService {
	return &RegistrationService{
		logger:            logger,
		js:                js,
		defaultUserPoolID: defaultUserPoolID,
		tenantsStream:     tenantsStream,
	}
}

// PublishTenantMessages publishes messages in response to a new tenant being onboarded.
func (rs *RegistrationService) PublishTenantMessages(ctx context.Context, id string, tenant model.NewTenant) error {
	// construct tenant message
	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}

	event := newTenantRegisteredMessage(values, id, tenant, rs.defaultUserPoolID)
	subject := fmt.Sprintf("%s.registered", rs.tenantsStream)
	bytes, err := event.Marshal()
	if err != nil {
		return nil
	}

	rs.js.Publish(subject, bytes)

	if Plan(tenant.Plan) == PlanPremium {
		if err = rs.provision(ctx); err != nil {
			return err
		}
	}

	return err
}

func (rs *RegistrationService) provision(_ context.Context) error {
	// start aws codepipeline
	return nil
}

func newTenantRegisteredMessage(values *web.Values, id string, tenant model.NewTenant, userPoolID string) msg.TenantRegisteredEvent {
	return msg.TenantRegisteredEvent{
		Metadata: msg.Metadata{
			TraceID: values.Metadata.TraceID,
			UserID:  values.Metadata.UserID,
		},
		Type: msg.TypeTenantRegistered,
		Data: msg.TenantRegisteredEventData{
			ID:         id,
			FullName:   tenant.FullName,
			Company:    tenant.Company,
			Email:      tenant.Email,
			Plan:       tenant.Plan,
			UserPoolID: userPoolID,
		},
	}
}

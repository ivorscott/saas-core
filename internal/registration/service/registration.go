package service

import (
	"context"
	"fmt"

	"github.com/devpies/saas-core/internal/registration/model"
	"github.com/devpies/saas-core/pkg/msg"

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
	appClientID       string
	defaultUserPoolID string
}

type Plan string

const (
	PlanBasic   Plan = "basic"
	PlanPremium Plan = "premium"
)

// NewRegistrationService returns a new registration service.
func NewRegistrationService(logger *zap.Logger, js publisher, appClientID string, defaultUserPoolID string) *RegistrationService {
	return &RegistrationService{
		logger:            logger,
		js:                js,
		appClientID:       appClientID,
		defaultUserPoolID: defaultUserPoolID,
	}
}

// PublishTenantMessages publishes messages in response to a new tenant being onboarded.
func (rs *RegistrationService) PublishTenantMessages(ctx context.Context, id string, tenant model.NewTenant) error {
	// construct tenant message
	event := newTenantRegisteredMessage(id, tenant)
	subject := fmt.Sprintf("%s.registered", rs.tenantsStream)
	bytes, err := event.Marshal()
	if err != nil {
		return nil
	}

	// publishRegistration
	rs.js.Publish(subject, bytes)

	// provision environment if necessary (	exit if not premium plan)
	if Plan(tenant.Plan) == PlanPremium {
		err = rs.provision(ctx)
		if err != nil {
			return err
		}
		// then, publish tenant config command
		command := newCreateConfigMessage(tenant, rs.appClientID, rs.defaultUserPoolID)
		subject = fmt.Sprintf("%s.configure", rs.tenantsStream)
		bytes, err = command.Marshal()
		if err != nil {
			return nil
		}

		// publishRegistration
		rs.js.Publish(subject, bytes)
	}

	return err
}

func (rs *RegistrationService) provision(ctx context.Context) error {
	// start aws codepipeline
	return nil
}

func newTenantRegisteredMessage(id string, tenant model.NewTenant) msg.TenantRegisteredEvent {
	return msg.TenantRegisteredEvent{
		Metadata: msg.Metadata{},
		Type:     msg.TypeTenantRegistered,
		Data: msg.TenantRegisteredEventData{
			ID:       id,
			FullName: tenant.FullName,
			Company:  tenant.Company,
			Email:    tenant.Email,
			Plan:     tenant.Plan,
		},
	}
}

func newCreateConfigMessage(tenant model.NewTenant, appClientID, userPoolID string) msg.CreateTenantConfigCommand {
	return msg.CreateTenantConfigCommand{
		Metadata: msg.Metadata{},
		Type:     msg.TypeCreateTenantConfig,
		Data: msg.CreateTenantConfigCommandData{
			TenantName:       tenant.Company,
			AppClientID:      appClientID,
			UserPoolID:       userPoolID,
			DeploymentStatus: "provisioned",
		},
	}
}

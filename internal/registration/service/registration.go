package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/devpies/saas-core/internal/registration/model"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"go.uber.org/zap"
)

// RegistrationService is responsible for managing tenant registration.
type RegistrationService struct {
	logger     *zap.Logger
	region     string
	idpService identityService
	js         publisher
}

type identityService interface {
	GetPlanBasedUserPool(ctx context.Context, tenant model.NewTenant, path string) (string, error)
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
func NewRegistrationService(logger *zap.Logger, region string, idpService identityService, js publisher) *RegistrationService {
	return &RegistrationService{
		logger:     logger,
		region:     region,
		idpService: idpService,
		js:         js,
	}
}

// CreateRegistration starts the tenant registration process.
func (rs *RegistrationService) CreateRegistration(ctx context.Context, tenantID string, tenant model.NewTenant) error {
	var err error
	userPoolID, err := rs.idpService.GetPlanBasedUserPool(ctx, tenant, rs.formatPath(tenant.Company))
	if err != nil {
		return err
	}
	err = rs.publishTenantRegisteredEvent(ctx, tenantID, tenant, userPoolID)
	if err != nil {
		return err
	}
	if err = rs.provision(ctx, Plan(tenant.Plan)); err != nil {
		return err
	}
	return nil
}

func (rs *RegistrationService) formatPath(company string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		rs.logger.Fatal("regex failed to compile", zap.Error(err))
	}
	return strings.ToLower(reg.ReplaceAllString(company, ""))
}

func (rs *RegistrationService) publishTenantRegisteredEvent(ctx context.Context, tenantID string, tenant model.NewTenant, userPoolID string) error {
	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}
	event := newTenantRegisteredEvent(values, tenantID, tenant, userPoolID)
	bytes, err := event.Marshal()
	if err != nil {
		return err
	}
	rs.js.Publish(msg.SubjectTenantRegistered, bytes)
	return nil
}

func (rs *RegistrationService) provision(ctx context.Context, plan Plan) error {
	if plan != PlanPremium {
		return nil
	}
	client := codepipeline.New(codepipeline.Options{
		Region: rs.region,
	})

	input := codepipeline.StartPipelineExecutionInput{
		Name:               aws.String("tenant-onboarding-pipeline"),
		ClientRequestToken: aws.String(fmt.Sprintf("requestToken-%s", time.Now().UTC())),
	}

	output, err := client.StartPipelineExecution(ctx, &input)
	if err != nil {
		return err
	}
	rs.logger.Info(fmt.Sprintf("successfully started pipeline - response: %+v", output))
	return nil
}

func newTenantRegisteredEvent(values *web.Values, tenantID string, tenant model.NewTenant, userPoolID string) msg.TenantRegisteredEvent {
	return msg.TenantRegisteredEvent{
		Metadata: msg.Metadata{
			TraceID: values.TraceID,
			UserID:  values.UserID,
		},
		Type: msg.TypeTenantRegistered,
		Data: msg.TenantRegisteredEventData{
			TenantID:   tenantID,
			FirstName:  tenant.FirstName,
			LastName:   tenant.LastName,
			Company:    tenant.Company,
			Email:      tenant.Email,
			Plan:       tenant.Plan,
			UserPoolID: userPoolID,
		},
	}
}

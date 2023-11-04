package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/devpies/saas-core/pkg/web"
	"net/http"
	"strings"
	"time"

	"github.com/devpies/saas-core/internal/tenant/model"
	"github.com/devpies/saas-core/pkg/msg"

	"go.uber.org/zap"
)

type publisher interface {
	Publish(subject string, message []byte)
}

type cognitoClient interface {
	AdminCreateUser(
		ctx context.Context,
		params *cognitoidentityprovider.AdminCreateUserInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

// TenantService manages tenant business operations.
type TenantService struct {
	logger         *zap.Logger
	js             publisher
	sharedPoolID   string
	cognitoClient  cognitoClient
	tenantRepo     tenantRepository
	connectionRepo connectionRepository
}

type tenantRepository interface {
	Insert(ctx context.Context, tenant model.NewTenant) error
	SelectOne(ctx context.Context, tenantID string) (model.Tenant, error)
	SelectAll(ctx context.Context) ([]model.Tenant, error)
	Update(ctx context.Context, id string, tenant model.UpdateTenant) error
	Delete(ctx context.Context, tenantID string) error
}

type connectionRepository interface {
	Insert(ctx context.Context, nc model.NewConnection) error
}

var (
	// ErrUsernameExistsException represents an AWS Cognito error caused by attempting to create and existing user.
	ErrUsernameExistsException = errors.New("account already exists")
)

// NewTenantService returns a new TenantService.
func NewTenantService(
	logger *zap.Logger,
	js publisher,
	sharedPoolID string,
	cognitoClient cognitoClient,
	tenantRepo tenantRepository,
	connectionRepo connectionRepository,
) *TenantService {
	return &TenantService{
		logger:         logger,
		js:             js,
		sharedPoolID:   sharedPoolID,
		cognitoClient:  cognitoClient,
		tenantRepo:     tenantRepo,
		connectionRepo: connectionRepo,
	}
}

// CreateTenantFromEvent creates a tenant from a message.
func (ts *TenantService) CreateTenantFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalTenantRegisteredEvent(m)
	if err != nil {
		return err
	}

	output, err := ts.createTenantIdentity(ctx, event.Data)
	if err != nil {
		ts.logger.Error("error creating cognito identity", zap.Error(err))
		switch {
		case strings.Contains(strings.ToLower(err.Error()), ErrUsernameExistsException.Error()):
			return web.NewRequestError(ErrUsernameExistsException, http.StatusBadRequest)
		default:
			return web.NewRequestError(err, http.StatusUnauthorized)
		}
	}
	userID := getUserID(output.User)
	created := output.User.UserCreateDate

	if err = ts.createTenant(ctx, event.Data, created); err != nil {
		ts.logger.Error("error storing tenant", zap.Error(err))
		return err
	}

	if err = ts.createTenantConnection(ctx, userID, event.Data.TenantID); err != nil {
		ts.logger.Error("error creating connection", zap.Error(err))
		return err
	}

	if err = ts.publish(ctx, userID, created, event.Data); err != nil {
		ts.logger.Error("error publishing message", zap.Error(err))
		return err
	}
	return nil
}

// FindOne finds a single tenant.
func (ts *TenantService) FindOne(ctx context.Context, tenantID string) (model.Tenant, error) {
	var (
		tenant model.Tenant
		err    error
	)
	tenant, err = ts.tenantRepo.SelectOne(ctx, tenantID)
	if err != nil {
		return tenant, err
	}
	return tenant, nil
}

// FindAll finds all tenants.
func (ts *TenantService) FindAll(ctx context.Context) ([]model.Tenant, error) {
	var err error
	tenants, err := ts.tenantRepo.SelectAll(ctx)
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

// Update updates a single tenant.
func (ts *TenantService) Update(ctx context.Context, id string, tenant model.UpdateTenant) error {
	var err error
	err = ts.tenantRepo.Update(ctx, id, tenant)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a tenant.
func (ts *TenantService) Delete(ctx context.Context, tenantID string) error {
	var err error
	err = ts.tenantRepo.Delete(ctx, tenantID)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TenantService) createTenantIdentity(ctx context.Context, data msg.TenantRegisteredEventData) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	return ts.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(ts.sharedPoolID),
		Username:   aws.String(data.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:tenant-id"), Value: aws.String(data.TenantID)},
			{Name: aws.String("custom:account-owner"), Value: aws.String("1")},
			{Name: aws.String("custom:company-name"), Value: aws.String(data.Company)},
			{Name: aws.String("custom:full-name"), Value: aws.String(fmt.Sprintf("%s %s", data.FirstName, data.LastName))},
			{Name: aws.String("email"), Value: aws.String(data.Email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	})
}

func (ts *TenantService) createTenant(ctx context.Context, data msg.TenantRegisteredEventData, created *time.Time) error {
	tenant := newTenant(data, created)
	return ts.tenantRepo.Insert(ctx, tenant)
}

func (ts *TenantService) createTenantConnection(ctx context.Context, userID, tenantID string) error {
	return ts.connectionRepo.Insert(ctx, model.NewConnection{
		UserID:   userID,
		TenantID: tenantID,
	})
}

func newTenant(data msg.TenantRegisteredEventData, created *time.Time) model.NewTenant {
	initialStatus := string(types.UserStatusTypeForceChangePassword)
	return model.NewTenant{
		ID:          data.TenantID,
		Email:       data.Email,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		CompanyName: data.Company,
		Plan:        data.Plan,
		Status:      initialStatus,
		Created:     *created,
	}
}

func newIdentityCreatedEvent(
	ctx context.Context,
	userID string,
	created *time.Time,
	data msg.TenantRegisteredEventData,
) (msg.TenantIdentityCreatedEvent, error) {
	var event msg.TenantIdentityCreatedEvent

	values, ok := web.FromContext(ctx)
	if !ok {
		return event, web.CtxErr()
	}

	event = msg.TenantIdentityCreatedEvent{
		Type: msg.TypeTenantIdentityCreated,
		Data: msg.TenantIdentityCreatedEventData{
			TenantID:  data.TenantID,
			UserID:    userID,
			Company:   data.Company,
			Email:     data.Email,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Plan:      data.Plan,
			CreatedAt: created.UTC().String(),
		},
		Metadata: msg.Metadata{
			TraceID: values.TraceID,
			UserID:  values.UserID,
		},
	}
	return event, nil
}

func getUserID(user *types.UserType) string {
	for _, v := range user.Attributes {
		if v.Name != nil && *v.Name == "sub" {
			return *v.Value
		}
	}
	return ""
}

func (ts *TenantService) publish(ctx context.Context, userID string, created *time.Time, data msg.TenantRegisteredEventData) error {
	e, err := newIdentityCreatedEvent(ctx, userID, created, data)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	ts.js.Publish(msg.SubjectTenantIdentityCreated, bytes)
	return nil
}

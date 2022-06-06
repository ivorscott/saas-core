package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"

	"go.uber.org/zap"
)

// RegistrationService is responsible for triggering tenant registration.
type RegistrationService struct {
	logger         *zap.Logger
	serviceAddress string
	servicePort    string
}

// NewRegistrationService returns a new registration service.
func NewRegistrationService(logger *zap.Logger, serviceAddress string, servicePort string) *RegistrationService {
	return &RegistrationService{
		logger:         logger,
		serviceAddress: serviceAddress,
		servicePort:    servicePort,
	}
}

// RegisterTenant sends new tenant to tenant registration microservice.
func (rs *RegistrationService) RegisterTenant(ctx context.Context, tenant model.NewTenant) error {
	data, err := json.Marshal(tenant)
	if err != nil {
		return err
	}

	payload := bytes.NewReader(data)
	url := fmt.Sprintf("%s:%s/register", rs.serviceAddress, rs.servicePort)

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		rs.logger.Info("registration failed", zap.Error(err))
	}
	defer resp.Body.Close()

	return err
}

//go:generate mockery --quiet --all --dir . --case snake --output ../mocks --exported
package handler_test

import (
	"github.com/devpies/saas-core/internal/admin/render"

	"github.com/stretchr/testify/mock"
)

var (
	testEmail          = "test@email.com"
	testPassword       = "password"
	testTemplateData   *render.TemplateData
	testIDToken        = "test_aws_cognito_id_token"
	testPChalSession   = "test_password_challenge_session"
	mockCtx            = mock.AnythingOfType("*context.valueCtx")
	fieldValidationErr = "field validation error"
)

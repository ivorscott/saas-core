package webapp

import (
	"context"

	"github.com/devpies/core/internal/admin/config"
	"github.com/devpies/core/internal/admin/render"

	"github.com/alexedwards/scs/v2"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"go.uber.org/zap"
)

type authService interface {
	Authenticate(ctx context.Context, email, password string) (*cip.AdminInitiateAuthOutput, error)
	CreateUserSession(ctx context.Context, token []byte) error
	RespondToNewPasswordRequiredChallenge(ctx context.Context, email, password string, session string) (*cip.AdminRespondToAuthChallengeOutput, error)
}

type WebApp struct {
	logger  *zap.Logger
	config  config.Config
	render  *render.Render
	service authService
	session *scs.SessionManager
}

func New(logger *zap.Logger, config config.Config, renderEngine *render.Render, service authService, session *scs.SessionManager) *WebApp {
	return &WebApp{
		logger:  logger,
		config:  config,
		render:  renderEngine,
		service: service,
		session: session,
	}
}

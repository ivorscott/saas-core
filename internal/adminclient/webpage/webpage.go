package webpage

import (
	"github.com/alexedwards/scs/v2"
	"github.com/devpies/core/internal/adminclient/config"
	"github.com/devpies/core/internal/adminclient/render"
	"go.uber.org/zap"
)

type WebPage struct {
	logger  *zap.Logger
	config  config.Config
	render  *render.Render
	session *scs.SessionManager
}

func New(logger *zap.Logger, config config.Config, renderEngine *render.Render, session *scs.SessionManager) *WebPage {
	return &WebPage{
		logger:  logger,
		config:  config,
		render:  renderEngine,
		session: session,
	}
}

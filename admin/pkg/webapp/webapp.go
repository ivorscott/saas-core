package webapp

import (
	"github.com/devpies/core/admin/pkg/config"
	"github.com/devpies/core/admin/pkg/webapp/render"
	"go.uber.org/zap"
)

type WebApp struct {
	logger *zap.Logger
	config config.Config
	render *render.Render
}

func New(logger *zap.Logger, config config.Config, renderEngine *render.Render) *WebApp {
	return &WebApp{
		logger: logger,
		config: config,
		render: renderEngine,
	}
}

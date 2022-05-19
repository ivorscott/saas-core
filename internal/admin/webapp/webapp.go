package webapp

import (
	"github.com/devpies/core/internal/admin"
	"github.com/devpies/core/internal/admin/webapp/render"
	"go.uber.org/zap"
)

type WebApp struct {
	logger *zap.Logger
	config admin.Config
	render *render.Render
}

func New(logger *zap.Logger, config admin.Config, renderEngine *render.Render) *WebApp {
	return &WebApp{
		logger: logger,
		config: config,
		render: renderEngine,
	}
}

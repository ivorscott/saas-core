package webpage

import (
	"github.com/devpies/core/internal/admin"
	"github.com/devpies/core/internal/admin/webpage/render"
	"go.uber.org/zap"
)

type WebPage struct {
	logger *zap.Logger
	config admin.ClientConfig
	render *render.Render
}

func New(logger *zap.Logger, config admin.ClientConfig, renderEngine *render.Render) *WebPage {
	return &WebPage{
		logger: logger,
		config: config,
		render: renderEngine,
	}
}

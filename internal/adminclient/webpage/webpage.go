package webpage

import (
	"github.com/devpies/core/internal/adminclient/config"
	"github.com/devpies/core/internal/adminclient/render"
	"go.uber.org/zap"
)

type WebPage struct {
	logger *zap.Logger
	config config.Config
	render *render.Render
}

func New(logger *zap.Logger, config config.Config, renderEngine *render.Render) *WebPage {
	return &WebPage{
		logger: logger,
		config: config,
		render: renderEngine,
	}
}

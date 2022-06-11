package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/tenant/model"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

type authInfoService interface {
	GetInfo(ctx context.Context, referer string) (model.AuthInfoAndRegion, error)
}

// AuthInfoHandler handles tenant requests.
type AuthInfoHandler struct {
	logger  *zap.Logger
	service authInfoService
}

// NewAuthInfoHandler returns a new AuthInfoHandler.
func NewAuthInfoHandler(logger *zap.Logger, service authInfoService) *AuthInfoHandler {
	return &AuthInfoHandler{
		logger:  logger,
		service: service,
	}
}

// GetAuthInfo handles a request for tenant authentication information.
func (ah *AuthInfoHandler) GetAuthInfo(w http.ResponseWriter, r *http.Request) error {
	referer := r.Header.Get("Referer")
	info, err := ah.service.GetInfo(r.Context(), referer)
	if err != nil {
		return web.NewRequestError(err, http.StatusNotFound)
	}
	return web.Respond(r.Context(), w, info, http.StatusOK)
}

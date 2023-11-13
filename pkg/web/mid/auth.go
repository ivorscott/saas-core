package mid

import (
	"net/http"

	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

// Auth middleware verifies the id_token.
func Auth(log *zap.Logger, region string, userPoolID string) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(handler web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			r, err := web.Authenticate(log, r, region, userPoolID)
			if err != nil {
				log.Info("api authentication failed", zap.Error(err))
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			return handler(w, r)
		}
		return h
	}
	return f
}

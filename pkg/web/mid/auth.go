package mid

import (
	"net/http"

	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

// Auth middleware verifies the id_token.
func Auth(log *zap.Logger, region string, userPoolClientID string) web.Middleware {
	f := func(after web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			r, err := web.Authenticate(log, r, region, userPoolClientID)
			if err != nil {
				log.Info("api authentication failed", zap.Error(err))
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			return after(w, r)
		}
		return h
	}
	return f
}

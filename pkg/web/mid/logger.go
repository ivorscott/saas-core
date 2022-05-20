package mid

import (
	"errors"
	"fmt"
	"github.com/devpies/core/pkg/web"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Logger writes some information about the request to the logs in the
// format: (200) GET /foo -> IP ADDR (latency)
func Logger(log *zap.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web value missing from context")
			}

			err := before(w, r)

			log.Info(fmt.Sprintf("(%d) : %s %s -> %s (%s)",
				v.StatusCode,
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				time.Since(v.Start)))

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return f
}

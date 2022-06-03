package mid

import (
	"github.com/devpies/core/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {

			// Run the handler chain and catch any propagated error.
			if err := before(w, r); err != nil {

				// Log the error.
				log.Info("", zap.Error(err))

				// Respond to the error.
				if err = web.RespondError(r.Context(), w, err); err != nil {
					return err
				}

				// Errors middleware is designed to consume errors but, if we receive a
				// shutdown error we need to return it so we can shutdown.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			// Return nil to indicate the error has been handled.
			return nil
		}

		return h
	}

	return f
}

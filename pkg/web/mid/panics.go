package mid

import (
	"fmt"
	"github.com/devpies/saas-core/pkg/web"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics(log *zap.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the error return
			// variable after the fact.
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic: %v", r)

					// Log the Go stack trace for this panic'd goroutine.
					log.Info(fmt.Sprintf("%s", debug.Stack()))
				}
			}()

			// Call the next Handler and set its return value in the error variable.
			return after(w, r)
		}

		return h
	}

	return f
}

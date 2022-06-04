// Package web provides a custom web framework.
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// App represents a new application.
type App struct {
	log      *zap.Logger
	mux      *chi.Mux
	mw       []Middleware
	shutdown chan os.Signal
}

// Handler represents a custom http handler that returns an error.
type Handler func(http.ResponseWriter, *http.Request) error

// Middleware runs some code before and/or after another Handler.
type Middleware func(Handler) Handler

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

// NewApp returns a new Web framework equipped with built-in middleware required for every handler.
func NewApp(router *chi.Mux, shutdown chan os.Signal, logger *zap.Logger, middleware ...Middleware) *App {
	return &App{
		log:      logger,
		mux:      router,
		mw:       middleware,
		shutdown: shutdown,
	}
}

// wrapMiddleware creates new handler by wrapping middleware around a handler.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// Loop backwards through the middleware list invoking each one.
	// Replaces the handler with the new wrapped handler. Looping backwards
	// ensures that the first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}
	return handler
}

// Handle converts our custom handler to the standard library Handler.
func (app *App) Handle(method string, url string, h Handler) {
	h = wrapMiddleware(app.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			Start: time.Now(),
		}

		ctx := r.Context()
		// Create a new context with new key/value.
		ctx = context.WithValue(ctx, KeyValues, &v)
		r = r.WithContext(ctx)
		// Catch any propagated error.
		if err := h(w, r); err != nil {
			app.log.Error("error: unhandled error", zap.Error(err))
			if IsShutdown(err) {
				app.SignalShutdown()
			}
		}
	}

	app.mux.MethodFunc(method, url, fn)
}

// ServeHTTP extends original mux ServeHTTP method.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}

// SignalShutdown sends application shutdown signal.
func (app *App) SignalShutdown() {
	app.log.Error("integrity issue: shutting down service")
	app.shutdown <- syscall.SIGSTOP
}

package web

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Values carries information about each request.
type Values struct {
	Metadata   *Metadata
	StatusCode int
	Start      time.Time
}

// Metadata is the type of value stored in the Context.
type Metadata struct {
	TraceID string
	UserID  string
	Token   string
}

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// NewContext returns a new Context that carries value r.
func NewContext(ctx context.Context, r *Values) context.Context {
	return context.WithValue(ctx, KeyValues, r)
}

// FromContext returns the Values value stored in ctx, if any.
func FromContext(ctx context.Context) (*Values, bool) {
	r, ok := ctx.Value(KeyValues).(*Values)
	return r, ok
}

// addContextMetadata adds Metadata to context.
func addContextMetadata(r *http.Request, token string, sub string) *http.Request {
	traceID := r.Header.Get("TraceID")
	if traceID == "" {
		traceID = uuid.New().String()
	}

	if v, ok := FromContext(r.Context()); ok {
		v.Metadata = &Metadata{
			TraceID: traceID,
			UserID:  sub,
			Token:   token,
		}
		ctx := NewContext(r.Context(), v)
		r = r.WithContext(ctx)
	}

	return r
}

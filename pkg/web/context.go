package web

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Values carries information about each request.
type Values struct {
	TenantID   string
	UserID     string
	Token      string
	TraceID    string
	StatusCode int
	Start      time.Time
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
func addContextMetadata(r *http.Request, token string, sub string, defaultTenantID string, tenantMap TenantConnectionMap) *http.Request {
	basePath := r.Header.Get("BasePath")
	traceID := r.Header.Get("TraceID")
	if traceID == "" {
		traceID = uuid.New().String()
	}

	if v, ok := FromContext(r.Context()); ok {
		v.UserID = sub
		v.Token = token
		v.TraceID = traceID
		v.TenantID = defaultTenantID

		val, okay := tenantMap[basePath]
		if okay {
			v.TenantID = val.TenantID
		}

		ctx := NewContext(r.Context(), v)
		r = r.WithContext(ctx)
	}

	return r
}

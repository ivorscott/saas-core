package web

// Middleware runs some code before and/or after another Handler.
type Middleware func(Handler) Handler

// wrapMiddleware creates new handler by wrapping middleware around a handler.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// Loop backwards through the middleware list invoking each one.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}
	return handler
}

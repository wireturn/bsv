package node

// A Middleware is a type that wraps a handler to remove boilerplate or other
// concerns not direct to any given Handler.
type Middleware func(Handler) Handler

// wrapMiddleware wraps a handler with some middleware.
func wrapMiddleware(handler Handler, mw []Middleware) Handler {

	// Wrap with our group specific middleware.
	for i := len(mw) - 1; i >= 0; i-- {
		if mw[i] != nil {
			handler = mw[i](handler)
		}
	}

	return handler
}

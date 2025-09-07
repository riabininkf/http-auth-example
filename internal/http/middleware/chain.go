package middleware

import "net/http"

// Middleware is a function that wraps an http.Handler to provide additional behavior for HTTP requests and responses.
type Middleware func(http.Handler) http.Handler

// Chain applies a stack of Middleware to an http.Handler, starting from the last middleware in the given slice.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

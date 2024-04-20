/*
Package server provides functionality for creating and managing HTTP servers, including middleware support.

Author: Your Name
*/

package server

import "net/http"

// Middleware defines the type for middleware functions that wrap http.Handler.
type Middleware func(http.Handler) http.HandlerFunc

// MiddlewareN defines the type for middleware functions that accept a ContextHandler.
type MiddlewareN func(handler ContextHandler)

// MiddlewareConvertedN defines the type for converted middleware functions that wrap http.Handler.
type MiddlewareConvertedN func(http.Handler) http.Handler

// MiddlewareChain creates a chain of middleware functions.
func (api *MyAPIServer) MiddlewareChain(middlewares []Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

// AddMiddleware adds middleware to the server's middleware list.
func (api *MyAPIServer) AddMiddleware(middleware Middleware) {
	if api.HandlerNew {
		api.Logger.Fatalf("You need to call this function AddMiddlewareN as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	api.Serv.MiddlewareList = append(api.Serv.MiddlewareList, middleware)
}

// AddMiddlewareN adds converted middleware to the server's middleware list.
func (api *MyAPIServer) AddMiddlewareN(middleware MiddlewareN) {
	if !api.HandlerNew {
		api.Logger.Fatalf("You need to call this function AddMiddleware as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	middlewareCN := MiddlewareWrapperN(middleware)
	api.Serv.MiddlewareListN = append(api.Serv.MiddlewareListN, middlewareCN)
}

// MiddlewareWrapperN wraps a middleware function that accepts a ContextHandler.
func MiddlewareWrapperN(handler func(ctx ContextHandler)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a ContextHandler with the provided ResponseWriter and Request
			ctx := ContextHandler{
				Writer:  w,
				Request: r,
			}

			// Call the middleware function with the ContextHandler
			handler(ctx)

			// Call the next handler with the ContextHandler
			next.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
}

// MiddlewareChainN creates a chain of converted middleware functions.
func (api *MyAPIServer) MiddlewareChainN(middleware []MiddlewareConvertedN) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		// Start with the innermost handler
		finalHandler := handler
		// Iterate over each middleware function in reverse order
		for i := len(middleware) - 1; i >= 0; i-- {
			// Apply the current middleware to the final handler
			finalHandler = middleware[i](finalHandler)
		}
		return finalHandler
	}
}

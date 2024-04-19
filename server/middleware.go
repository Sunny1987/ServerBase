package server

import "net/http"

type Middleware func(http.Handler) http.HandlerFunc

func (api *MyAPIServer) MiddlewareChain(middlewares []Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}

}

func (api *MyAPIServer) AddMiddleware(middleware Middleware) {
	api.Serv.MiddlewareList = append(api.Serv.MiddlewareList, middleware)
}

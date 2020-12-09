package httpxroute

import "github.com/gin-gonic/gin"

type (
	Router struct {
		Method string
		Path   string
		fn     interface{}
	}
)

func NewHandlerFunc(fn interface{}, opts ...HandlerOption) gin.HandlerFunc {
	return CreateHandlerFunc(fn, opts...)
}

func Route(routes gin.IRoutes, method string, path string, fn interface{}, opts ...HandlerOption) gin.IRoutes {
	return routes.Handle(method, path, NewHandlerFunc(fn, opts...))
}

func RouteWithRouter(routes gin.IRoutes, r Router, opts ...HandlerOption) gin.IRoutes {
	return Route(routes, r.Method, r.Path, r.fn, opts...)
}

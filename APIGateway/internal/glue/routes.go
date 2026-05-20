package glue

import (
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Method      string
	Path        string
	Handler     func(*gin.Context)
	Middlewares []func(logger *platform.Logger) gin.HandlerFunc
}

func RegisterRoutes(router *gin.RouterGroup, routes []Route,logger *platform.Logger) {
	for _, route := range routes {

		var handlers []gin.HandlerFunc
		for _, mw := range route.Middlewares {
			
			
			handlers = append(handlers, mw(logger))
		}

		
		handlers = append(handlers, route.Handler)
		router.Handle(route.Method, route.Path, handlers...)
	}
}
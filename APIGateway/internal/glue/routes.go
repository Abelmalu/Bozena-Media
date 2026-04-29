package glue

import "github.com/gin-gonic/gin"

type Route struct {
	Method      string
	Path        string
	Handler     func(*gin.Context)
	Middlewares []func() gin.HandlerFunc
}

func RegisterRoutes(router *gin.RouterGroup, routes []Route) {
	for _, route := range routes {

		var handlers []gin.HandlerFunc
		for _, mw := range route.Middlewares {
			
			
			handlers = append(handlers, mw())
		}

		
		handlers = append(handlers, route.Handler)
		router.Handle(route.Method, route.Path, handlers...)
	}
}
package routing

import (
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/glue"
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)
func InitAuthRoute(router *gin.RouterGroup, handler *handler.AuthHandler,logger *platform.Logger) {


	routes := []glue.Route{
		{
			Method:      http.MethodPost,
			Path:        "/register",
			Handler:     handler.Register,
			Middlewares: []func(logger *platform.Logger) gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/login",
			Handler:     handler.Login,
			Middlewares: []func(logger *platform.Logger) gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/refresh",
			Handler:     handler.RefreshHandler,
			Middlewares: []func(logger *platform.Logger) gin.HandlerFunc{},
		},
		{
			Method:  http.MethodPost,
			Path:    "/logout",
			Handler: handler.Logout,
			Middlewares: []func(logger *platform.Logger) gin.HandlerFunc{
				middleware.AuthMiddleware,
			},
		},
	}

	glue.RegisterRoutes(router,routes,logger)

}

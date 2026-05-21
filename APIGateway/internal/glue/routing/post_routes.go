package routing

import (
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/glue"
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)
func InitPostRoute(router *gin.RouterGroup, handler *handler.PostHandler,logger *platform.Logger) {

	routes := []glue.Route{
		{
			Method:      http.MethodPost,
			Path:        "/",
			Handler:     handler.CreatePost,
			Middlewares: []func(*platform.Logger) gin.HandlerFunc{

				middleware.AuthMiddleware,

			
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "/",
			Handler:     handler.ListPosts,
			Middlewares: []func(*platform.Logger) gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/update/:id",
			Handler:     handler.UpdatePost,
			Middlewares: []func(*platform.Logger) gin.HandlerFunc{
				middleware.AuthMiddleware,

			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/delete/:id",
			Handler: handler.DeletePost,
			Middlewares: []func(*platform.Logger) gin.HandlerFunc{
				middleware.AuthMiddleware,
			},
		},
	}

	glue.RegisterRoutes(router,routes,logger)

}

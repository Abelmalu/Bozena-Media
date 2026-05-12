package routing

import (
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/glue"
    "github.com/abelmalu/golang-posts/APIGateway/internal/handler"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/gin-gonic/gin"
)
func InitPostRoute(router *gin.RouterGroup, handler *handler.PostHandler) {

	routes := []glue.Route{
		{
			Method:      http.MethodPost,
			Path:        "/",
			Handler:     handler.CreatePost,
			Middlewares: [
				
			]func() gin.HandlerFunc{middleware.AuthMiddleware},
		},
		{
			Method:      http.MethodGet,
			Path:        "/",
			Handler:     handler.ListPosts,
			Middlewares: []func() gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/update/:id",
			Handler:     handler.UpdatePost,
			Middlewares: []func() gin.HandlerFunc{},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/delete/:id",
			Handler: handler.DeletePost,
			Middlewares: []func() gin.HandlerFunc{
				middleware.AuthMiddleware,
			},
		},
	}

	glue.RegisterRoutes(router,routes)

}

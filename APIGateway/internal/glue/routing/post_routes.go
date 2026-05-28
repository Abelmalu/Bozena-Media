package routing

import (
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/glue"
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)
func InitPostRoute(router *gin.RouterGroup, handler *handler.PostHandler,logger *platform.Logger,redisClient *redis.Client) {

	routes := []glue.Route{
		{
			Method:      http.MethodPost,
			Path:        "/",
			Handler:     handler.CreatePost,
			Middlewares: [] gin.HandlerFunc{

				middleware.AuthMiddleware(logger,redisClient),

			
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "/",
			Handler:     handler.ListPosts,
			Middlewares: [] gin.HandlerFunc{

				middleware.AuthMiddleware(logger,redisClient),
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "/user/:id",
			Handler:     handler.GetUserPosts,
			Middlewares: [] gin.HandlerFunc{

				middleware.AuthMiddleware(logger,redisClient),
			},
		},
		{
			Method:      http.MethodPut,
			Path:        "/update/:id",
			Handler:     handler.UpdatePost,
			Middlewares: []gin.HandlerFunc{
				middleware.AuthMiddleware(logger,redisClient),

			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/delete/:id",
			Handler: handler.DeletePost,
			Middlewares: [] gin.HandlerFunc{
				middleware.AuthMiddleware(logger,redisClient),
			},
		},
	}

	glue.RegisterRoutes(router,routes)

}

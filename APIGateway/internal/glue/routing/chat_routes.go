package routing

import (
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/glue"
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
//	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitChatRoute(router *gin.RouterGroup, handler *handler.ChatHandler, logger *platform.Logger, redis *redis.Client) {

	routes := []glue.Route{

		{
			Method:  http.MethodGet,
			Path:    "/ws",
			Handler: handler.Connect,
			Middlewares: []gin.HandlerFunc{

				//middleware.AuthMiddleware(logger, redis),
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}

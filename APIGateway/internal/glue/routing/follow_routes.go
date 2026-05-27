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



func InitFollowRoute(router *gin.RouterGroup, handler *handler.FollowHandler,logger *platform.Logger,redisClient *redis.Client){

	routes := []glue.Route{

		{
			Method: http.MethodPost,
			Path:"/:id",
			Handler: handler.ToggleFollow,
			Middlewares: [] gin.HandlerFunc{

				middleware.AuthMiddleware(logger,redisClient),
			},


		},


	}
	glue.RegisterRoutes(router,routes)
}
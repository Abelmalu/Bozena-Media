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







func InitLikeRoute(router *gin.RouterGroup, handler *handler.LikeHandler,logger *platform.Logger,redisClient *redis.Client){

	routes := []glue.Route{

		{
			Method: http.MethodPost,
			Path:"/like/:id",
			Handler: handler.ToggleLike,
			Middlewares: [] gin.HandlerFunc{

				middleware.AuthMiddleware(logger,redisClient),
			},


		},
		{
			Method: http.MethodGet,
			Path:"/likes/:id",
			Handler: handler.GetPostLikes,
			Middlewares: [] gin.HandlerFunc{

				middleware.AuthMiddleware(logger,redisClient),
			},


		},


	}
	glue.RegisterRoutes(router,routes)
}
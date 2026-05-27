package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/internal/glue/routing"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitRoute(router *gin.Engine, handler Handler,logger *platform.Logger,redisClient *redis.Client) {

	// Global middlewares
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.RequestLoggerMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.RateLimitMiddleware(redisClient,logger))

	authRouter := router.Group("api/auth")
	routing.InitAuthRoute(authRouter,handler.authHandler,logger,redisClient)

	//post routes initialization
	postRouter := router.Group("api/posts")
	routing.InitPostRoute(postRouter,handler.postHandler,logger,redisClient)

	//like route initialization
	likeRouter := router.Group("api/posts")

    routing.InitLikeRoute(likeRouter,handler.likeHandler,logger,redisClient)



}
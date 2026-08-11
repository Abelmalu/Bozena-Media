package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/internal/glue/routing"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitRoute(router *gin.Engine, handler Handler, logger *platform.Logger, redisClient *redis.Client) {

	// Global middlewares
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.RequestLoggerMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.RateLimitMiddleware(redisClient, logger))
	router.Use(middleware.CORSMiddleware())

	authRouter := router.Group("api/auth")
	routing.InitAuthRoute(authRouter, handler.ah, logger, redisClient)

	//post routes initialization
	postRouter := router.Group("api/posts")
	routing.InitPostRoute(postRouter, handler.ps, logger, redisClient)

	//like route initialization
	likeRouter := router.Group("api/post")
	routing.InitLikeRoute(likeRouter, handler.lh, logger, redisClient)

	//follow route initialization
	followRouter := router.Group("api/follow")
	routing.InitFollowRoute(followRouter, handler.fh, logger, redisClient)

	//follow route initialization

	feedRouter := router.Group("api/feed")
	routing.InitFeedRoute(feedRouter, handler.fd, logger, redisClient)

	// notification routes
	notificationRouter := router.Group("api/notification")
	routing.InitNotificationRoutes(notificationRouter, handler.notificationHandler, logger, redisClient)


	// chat routes
	chatRouter := router.Group("api/chat")
	routing.InitChatRoute(chatRouter,handler.ch,logger,redisClient)
}

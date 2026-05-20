package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/internal/glue/routing"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.Engine, handler Handler,logger *platform.Logger) {

	// Global middlewares
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.RequestLoggerMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))

	authRouter := router.Group("api/auth")
	routing.InitAuthRoute(authRouter,&handler.authHandler,logger)

	//post routes initialization
	postRouter := router.Group("api/posts")
	routing.InitPostRoute(postRouter,&handler.postHandler,logger)




}
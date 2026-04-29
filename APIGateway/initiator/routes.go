package initiator


import (
	"github.com/abelmalu/golang-posts/APIGateway/internal/glue/routing"
	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.Engine,handler Handler){

	authRouter := router.Group("api/auth")
	routing.InitAuthRoute(authRouter,&handler.authHandler)

	//post routes initialization
	postRouter := router.Group("/api/posts")
	routing.InitPostRoute(postRouter,&handler.postHandler)




}
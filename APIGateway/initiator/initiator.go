package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/cmd/server"
	"github.com/gin-gonic/gin"
)

func Initialize() {

	router := gin.New()
	router.Use(server.CORSMiddleware())

	//initializing loggers 
	logger := InitLogger()

	//Initializing redis client
	redisClient := InitRedis("127.0.0.1:6379","",0,logger)



	// initializing clients
	clients := NewClient(logger)

	// initializing handlers
	handlers := InitHandler(*clients,logger)

	//initializing routes
	InitRoute(router, *handlers,logger,redisClient)

	//start the gin server 
	server.StartServer(router)


}

package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/cmd/server"
	"github.com/gin-gonic/gin"
)

func Initialize() {

	router := gin.New()

	//initializing loggers 
	logger := InitLogger()

	// initializing clients
	clients := NewClient(logger)

	// initializing handlers
	handlers := InitHandler(*clients,logger)

	//initializing routes
	InitRoute(router, *handlers,logger)

	//start the gin server 
	server.StartServer(router)


}

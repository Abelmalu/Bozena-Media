package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/cmd/server"
	"github.com/gin-gonic/gin"
)

func Initialize() {

	router := gin.Default()

	// initializing clients
	clients := NewClient()

	// initializing handlers
	handlers := InitHandler(*clients)

	//initializing routes
	InitRoute(router, *handlers)

	//start the gin server 
	server.StartServer(router)


}

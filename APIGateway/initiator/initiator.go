package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/cmd/server"
	"github.com/abelmalu/golang-posts/APIGateway/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Initialize() {

	router := gin.New()

	//initializing loggers
	logger := InitLogger()

	//load environment variables
	cfg, err := config.LoadConfig()

	if err != nil {

		logger.Error("Error while loading environment variables", zap.Error(err))

	}

	//Initializing redis client
	redisClient := InitRedis("127.0.0.1:6379", "", 0, logger)

	// initializing clients
	clients := NewClient(logger,cfg)

	// initializing handlers
	handlers := InitHandler(*clients, logger)

	//initializing routes
	InitRoute(router, *handlers, logger, redisClient)

	//start the gin server
	server.StartServer(router)

}

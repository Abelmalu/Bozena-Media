package initiator

import (
	"context"
	"log"

	"github.com/abelmalu/golang-posts/platform"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)


func InitRedis(address,password string,db int,logger *platform.Logger) (*redis.Client){
	ctx := context.Background()

	redisClient := redis.NewClient(
		&redis.Options{
			Addr: address,
			Password: password,
			DB: db,
		},
	)

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	logger.Info("Redis connected successfully!",zap.String("Response",pong))

	return redisClient


}
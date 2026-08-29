package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostServiceADD         string
	AuthServiceADD         string
	LikeServiceADD         string
	FeedServiceADD         string
	FollowServiceADD       string
	NotificationServiceADD string
	ChatServiceADD         string
	RedisADD               string
	RedisPassword          string
	RedisDB                int
	Port                   string
}

// loads configuration varibales from the environment and inject them to the config struct
func LoadConfig() (*Config, error) {

	config := Config{}


	config.AuthServiceADD = os.Getenv("AUTH_SERVICE_URL")
	if config.AuthServiceADD == "" {
		return nil, fmt.Errorf("AUTH_SERVICE_URL is required")
	}
	
	config.PostServiceADD = os.Getenv("POST_SERVICE_URL")
	if config.PostServiceADD == "" {
		return nil, fmt.Errorf("POST_SERVICE_URL is required")
	}



	config.LikeServiceADD = os.Getenv("LIKE_SERVICE_URL")
	if config.LikeServiceADD == "" {
		return nil, fmt.Errorf("LIKE_SERVICE_URL is required")
	}

	config.FeedServiceADD = os.Getenv("FEED_SERVICE_URL")
	if config.FeedServiceADD == "" {
		return nil, fmt.Errorf("FEED_SERVICE_URL is required")
	}

	config.FollowServiceADD = os.Getenv("FOLLOW_SERVICE_URL")
	if config.FollowServiceADD == "" {
		return nil, fmt.Errorf("FOLLOW_SERVICE_URL is required")
	}

	config.NotificationServiceADD = os.Getenv("NOTIFICATION_SERVICE_URL")
	if config.NotificationServiceADD == "" {
		return nil, fmt.Errorf("NOTIFICATION_SERVICE_URL is required")
	}

	config.ChatServiceADD = os.Getenv("CHAT_SERVICE_URL")
	if config.ChatServiceADD == "" {
		return nil, fmt.Errorf("CHAT_SERVICE_URL is required")
	}

	config.RedisADD = os.Getenv("REDIS_URL")
	if config.RedisADD == "" {
		return nil, fmt.Errorf("REDIS_ADD is required")
	}

	config.RedisPassword = os.Getenv("REDIS_PASSWORD") 

	
	config.Port = ":" + os.Getenv("SERVER_PORT")
	if config.Port == "" {

		config.Port = ":8080"
	}

	return &config, nil

}

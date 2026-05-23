package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostServiceADD string
	AuthServiceADD string
	LikeServiceADD string
	RedisADD string
	RedisPassword string
	RedisDB int
	Port           string
}

// loads configuration varibales from the environment and inject them to the config struct
func LoadConfig() (*Config, error) {

	config := Config{}

	config.PostServiceADD = os.Getenv("Post_Service_ADD")

	if config.PostServiceADD  == "" {

		return nil,fmt.Errorf("Post_Service_ADD is required")
	}

	config.AuthServiceADD = os.Getenv("Auth_Service_ADD")
	if config.AuthServiceADD == ""{

		return nil,fmt.Errorf("Auth_Service_ADD is required")
	}

	config.RedisADD = os.Getenv("REDIS_ADD")
	if config.RedisADD == ""{

		return nil,fmt.Errorf("Auth_Service_ADD is required")
	}




	config.Port = os.Getenv("PORT")
	if config.Port == ""{

		config.Port = "8080"
	}

	return &config,nil

}

package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)


type Config struct{
	DBURL string 
	GRPCPORT int
	KafkaBrokersURL []string
	RedisADD   string

}
func LoadConfig() (*Config,error){


	cfg := Config{}

	var err error

	portStr := os.Getenv("GRPC_PORT")

	if portStr == "" {

		portStr = "50054"
	}
	cfg.GRPCPORT, err = strconv.Atoi(portStr)
	if err != nil {

		return nil, fmt.Errorf("invalid SERVER_PORT '%s': must be an integer", portStr)

	}

	cfg.DBURL = os.Getenv("DB_URL")

	if cfg.DBURL == ""{


		return nil,errors.New("DB_URL environment variable is required!")
	}

	cfg.KafkaBrokersURL = []string {os.Getenv("KAFKA_BROKERS_URL")}
	if len(cfg.KafkaBrokersURL) == 0 {
		return nil, fmt.Errorf("MINIO_ADD is required")
	}

		cfg.RedisADD = os.Getenv("REDIS_URL")
	if cfg.RedisADD == "" {
		return nil, fmt.Errorf("REDIS_ADD is required")
	}


	return &cfg,nil

}



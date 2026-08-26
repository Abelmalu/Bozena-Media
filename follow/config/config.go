package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DBURL           string
	GRPCPORT        string
	KafkaBrokersURL []string
	RedisADD        string
}

func LoadConfig() (*Config, error) {

	cfg := Config{}


	cfg.GRPCPORT = ":" + os.Getenv("GRPC_PORT")

	if cfg.GRPCPORT == "" {

		return nil, errors.New("GRPC_POR environment variable is required!")

	}

	cfg.DBURL = os.Getenv("DB_URL")

	if cfg.DBURL == "" {

		return nil, errors.New("DB_URL environment variable is required!")
	}

	cfg.KafkaBrokersURL = []string{os.Getenv("KAFKA_BROKERS_URL")}
	if len(cfg.KafkaBrokersURL) == 0 {
		return nil, fmt.Errorf("MINIO_ADD is required")
	}

	cfg.RedisADD = os.Getenv("REDIS_URL")
	if cfg.RedisADD == "" {
		return nil, fmt.Errorf("REDIS_ADD is required")
	}

	return &cfg, nil

}

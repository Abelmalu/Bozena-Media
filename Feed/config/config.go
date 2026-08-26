package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DBURL            string
	GRPCPORT         string
	FollowServiceADD string
	MinioADD		 string
	KafkaBrokersURL         []string
}

func LoadConfig() (*Config, error) {

	cfg := Config{}

	cfg.GRPCPORT = ":" + os.Getenv("GRPC_PORT")
	if cfg.GRPCPORT == "" {
		cfg.GRPCPORT = ":50055"
	}

	cfg.DBURL = os.Getenv("DB_URL")
	if cfg.DBURL == "" {
		return nil, errors.New("DB_URL environment variable is required!")
	}

	cfg.FollowServiceADD = os.Getenv("FOLLOW_SERVICE_URL")
	if cfg.FollowServiceADD == "" {
		return nil, fmt.Errorf("FOLLOW_SERVICE_URL is required")
	}

	cfg.MinioADD = os.Getenv("MINIO_URL")
	if cfg.MinioADD == "" {
		return nil, fmt.Errorf("MINIO_ADD is required")
	}

	cfg.KafkaBrokersURL = []string {os.Getenv("KAFKA_BROKERS_URL")}
	if len(cfg.KafkaBrokersURL) == 0 {
		return nil, fmt.Errorf("MINIO_ADD is required")
	}

	return &cfg, nil

}

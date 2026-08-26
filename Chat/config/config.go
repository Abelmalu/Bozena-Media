package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DBURL    string
	MinioADD string
	Port     string
	KafkaBrokersURL  []string
}

func LoadConfig() (*Config, error) {

	cfg := Config{}

	cfg.DBURL = os.Getenv("DB_URL")

	if cfg.DBURL == "" {

		return nil, errors.New("DB_URL environment variable is required!")
	}

	cfg.MinioADD = os.Getenv("MINIO_URL")
	if cfg.MinioADD == "" {
		return nil, fmt.Errorf("MINIO_ADD is required")
	}

	cfg.Port = ":" + os.Getenv("SERVER_PORT")
	if cfg.Port == "" {

		cfg.Port = ":8084"
	}


	cfg.KafkaBrokersURL = []string {os.Getenv("KAFKA_BROKERS_URL")}
	if len(cfg.KafkaBrokersURL) == 0 {
		return nil, fmt.Errorf("MINIO_ADD is required")
	}

	return &cfg, nil

}

package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DBURL           string
	GRPCPORT        int
	ServerPort      string
	KafkaBrokersURL []string
}

func LoadConfig() (*Config, error) {

	cfg := Config{}

	cfg.ServerPort = ":" + os.Getenv("SERVER_PORT")

	if cfg.ServerPort == "" {

		return nil, errors.New("PORT environment variable is required!")

	}

	cfg.DBURL = os.Getenv("DB_URL")

	if cfg.DBURL == "" {

		return nil, errors.New("DB_URL environment variable is required!")
	}

	cfg.KafkaBrokersURL = []string{os.Getenv("KAFKA_BROKERS_URL")}
	if len(cfg.KafkaBrokersURL) == 0 {
		return nil, fmt.Errorf("KAKFKA address is required")
	}

	return &cfg, nil

}

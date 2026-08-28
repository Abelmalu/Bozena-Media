package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)


type Config struct {
	DBURL           string
	GRPCPORT        int
	KafkaBrokersURL []string
}

func LoadConfig() (*Config, error) {

	cfg := Config{}

	var err error

	portStr := os.Getenv("GRPC_PORT")

	if portStr == "" {

		portStr = "50053"
	}
	cfg.GRPCPORT, err = strconv.Atoi(portStr)
	if err != nil {

		return nil, fmt.Errorf("invalid SERVER_PORT '%s': must be an integer", portStr)

	}

	cfg.DBURL = os.Getenv("DB_URL")

	if cfg.DBURL == "" {

		return nil, errors.New("DB_URL environment variable is required!")
	}

	cfg.KafkaBrokersURL = []string{os.Getenv("KAFKA_BROKERS_URL")}
	if len(cfg.KafkaBrokersURL) == 0 {
		return nil, fmt.Errorf("KAFKA_PRODUCER is required")
	}
	return &cfg, nil

}



package config

import (
	"errors"
	"os"
)

type Config struct {
	DBURL  string
	DBName string
}

func LoadConfig() (*Config, error) {

	cfg := Config{}
	
	cfg.DBURL = os.Getenv("DB_URL")

	if cfg.DBURL == "" {

		return nil, errors.New("DB_URL environment variable is required!")
	}

	return &cfg, nil

}

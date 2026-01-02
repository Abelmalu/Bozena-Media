package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct{
	ServerPort int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	


}

func LoadConfig() (*Config, error) {

	cfg := &Config{}
	var err error

	// 1. Load Server Configuration
	portStr := os.Getenv("SERVER_PORT")

	if portStr == "" {

		portStr = "8080"
	}
	cfg.ServerPort, err = strconv.Atoi(portStr)
	if err != nil {

		return nil, fmt.Errorf("invalid SERVER_PORT '%s': must be an integer", portStr)

	}
	cfg.DBHost = os.Getenv("DB_HOST")
	if cfg.DBHost == "" {

		return nil, fmt.Errorf("DB_HOST environment variable is required")
	}
	dbPortStr := os.Getenv("DB_PORT")

	if dbPortStr == "" {

		return nil, fmt.Errorf("DB_PORT environment variable is required")
	}
	cfg.DBPort, err = strconv.Atoi(dbPortStr)

	if err != nil {

		return nil, fmt.Errorf("invalid SERVER_PORT '%s': must be an integer", dbPortStr)

	}

	cfg.DBName = os.Getenv("DB_NAME")

	// if cfg.DBName == "" {

	// 	return nil, fmt.Errorf("DB_NAME environment variable is required")
	// }
	cfg.DBUser = os.Getenv("DB_USER")
	if cfg.DBUser == "" {

		return nil, fmt.Errorf("DB_USER environment variable is required")

	}
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	// if cfg.DBPassword == "" {

	// 	return nil, fmt.Errorf("DB_Password environment variable is required")

	// }


	return cfg, nil

}
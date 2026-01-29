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

}
func LoadConfig() (*Config,error){


	cfg := Config{}

	var err error

	portStr := os.Getenv("GRPC_PORT")

	if portStr == "" {

		portStr = "8080"
	}
	cfg.GRPCPORT, err = strconv.Atoi(portStr)
	if err != nil {

		return nil, fmt.Errorf("invalid SERVER_PORT '%s': must be an integer", portStr)

	}

	cfg.DBURL = os.Getenv("DB_URL_POST")

	if cfg.DBURL == ""{


		return nil,errors.New("DB_URL environment variable is required!")
	}


	return &cfg,nil

}



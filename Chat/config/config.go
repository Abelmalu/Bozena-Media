package config

import (
	"errors"
	"os"
)


type Config struct{
	DBURL string 
	DBName string
	

}
func LoadConfig() (*Config,error){


	cfg := Config{}


	portStr := os.Getenv("GRPC_PORT")

	if portStr == "" {

		portStr = "50051"
	}


	cfg.DBName = os.Getenv("DB_NAME")

	if cfg.DBName == ""{

		return nil,errors.New("DB_NAME environment variable is required!")



	}

	cfg.DBURL = os.Getenv("DB_URL")

	if cfg.DBURL == ""{


		return nil,errors.New("DB_URL environment variable is required!")
	}


	return &cfg,nil

}



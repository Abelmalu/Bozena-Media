package main

import (
	"log"

	"github.com/abelmalu/golang-posts/APIGateway/initiator"
	"github.com/joho/godotenv"
)

func main() {

	// load environment variables using godoenv package
	if err := godotenv.Load(); err != nil {

		log.Fatalf("Error while loading environment variables %v", err)

	}

	initiator.Initialize()

}

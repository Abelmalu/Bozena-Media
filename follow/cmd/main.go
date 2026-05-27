package main

import (
	"log"

	 "github.com/abelmalu/golang-posts/follow/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables using godoenv package
	if err := godotenv.Load(); err != nil {

		log.Fatalf("Error while loading environment variables %v", err)

	}

	//Initiating the application
	app := application.NewApp()
	
	// Run the created application instance/grpc server
	app.Run()

}

package main

import (
	"log"

	application "github.com/abelmalu/golang-posts/Auth/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables using godoenv package
	if err := godotenv.Load(); err != nil {

		log.Fatalf("Error while loading environment variables %v", err)

	}

	//Initiating the application
	app, err := application.NewApp()
	if err != nil {

		log.Fatalf("Application Initializing Error %v", err)
	}
	// Run the created application instance/grpc server
	app.Run()

}

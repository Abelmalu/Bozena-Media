package main

import (
	"log"

	application "github.com/abelmalu/golang-posts/Chat/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	// load environment variables using godoenv package
	if err := godotenv.Load(); err != nil {

		log.Printf("Error while loading environment variables %v", err)

	}

	//Initiating the application
	app := application.NewApp()
	
	// Run the created application instance/grpc server
	app.Run()

}

package main

import (
	"log"

	"github.com/abelmalu/golang-posts/post/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	//Initiating the application
	if err := godotenv.Load(); err != nil{

		log.Fatalf("Error while loading environment variables %v",err)

	}

	app, err := application.NewApp()
	if err != nil {

		log.Fatalf("Application Initializing Error %v", err)
	}

	app.Run()

}

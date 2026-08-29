package main

import (
	"log"

	application "github.com/abelmalu/golang-posts/notification/internal/app"
	"github.com/joho/godotenv"
)


func main(){

	if err := godotenv.Load(); err != nil {

		log.Printf("Error loading env variables %v",err)
	}

	app := application.NewApp()

	app.Run()
}
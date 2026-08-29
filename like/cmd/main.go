package main

import (
	"log"

	application "github.com/abelmalu/golang-posts/like/internal/app"
	"github.com/joho/godotenv"
)


func main(){

// load environment variables using godoenv package
	if err := godotenv.Load(); err != nil {

		log.Printf("Error while loading environment variables %v", err)

	}

app := application.NewApp()

app.Run()


}

package main

import (
	"log"

	"github.com/abelmalu/golang-posts/APIGateway/internal/clients"
	"github.com/abelmalu/golang-posts/APIGateway/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	
)

func main() {


	if err := godotenv.Load(); err != nil{

		log.Fatalf("Error while loading env variabales %v",err)
	}
	postConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to gRPC server: %v", err)

	}

	defer postConn.Close()
    
	postClient := clients.NewPostClient(postConn)

	postHandler := handlers.NewPostHandler(postClient)

	r := gin.Default()

	r.POST("/posts", postHandler.CreatePost)

	r.Run(":8080")

}

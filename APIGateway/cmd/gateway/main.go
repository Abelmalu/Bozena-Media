package main

import (
	"log"

	"github.com/abelmalu/golang-posts/APIGateway/internal/clients"
	"github.com/abelmalu/golang-posts/APIGateway/internal/handlers"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"

	//"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
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
	authConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to gRPC server: %v", err)

	}

	defer postConn.Close()

    //Dependency Injection for PostService
	postClient := client.NewPostClient(postConn)
	postHandler := handler.NewPostHandler(postClient)

	//Dependency Injectiion for AuthService 
	authClient := client.NewAuthClient(authConn)
	authHandler := handler.NewAuthHandler(authClient)



	r := gin.Default()

	//Authentication route group
	auth := r.Group("api/auth")
	{
		{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		// auth.POST("/refresh", authHandler.RefreshHandler)
		// auth.POST("/logout", middleware.AuthMiddleware(), auth.Logout)

	}
	}

	// post route group
	post := r.Group("/api/posts")
	post.Use(middleware.AuthMiddleware())
	{
		post.POST("/", postHandler.CreatePost)

	}

	

	r.Run(":8080")

}

package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abelmalu/golang-posts/Chat/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type App struct {
	mongoClient *mongo.Client
}

func NewApp() *App {

	cfg,err :=config.LoadConfig()

	if err != nil {

		log.Fatalf("Error loading environment variables",err)
	}

	client := initDB(cfg)

	return &App{

		mongoClient:client ,
	}
}

func initDB(cfg *config.Config) *mongo.Client {


	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.DBURL))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	// 6. Access a specific database and collection
	collection := client.Database("mydatabase").Collection("mycollection")
	_ = collection


	return client
}

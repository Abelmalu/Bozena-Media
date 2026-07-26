package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abelmalu/golang-posts/Chat/config"
	"github.com/abelmalu/golang-posts/Chat/internal/handlers"
	"github.com/abelmalu/golang-posts/Chat/internal/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type App struct {
	mongoClient *mongo.Database
}

func NewApp() *App {

	cfg, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("Error loading environment variables %v", err)
	}

	DB := initDB(cfg)

	return &App{

		mongoClient: DB,
	}
}

func initDB(cfg *config.Config) *mongo.Database {

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

	database := client.Database("chat_db")

	return database
}

func InitRoute(h *handlers.ChatHandler) {

	router := gin.New()


	router.Handle(http.MethodGet,"/api/chat/ws",h.HandleWebSocket)
}

func (app *App) Run(){

	chatRepo := repository.NewChatRepository(app.mongoClient)



	





}

package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abelmalu/golang-posts/Chat/config"
	"github.com/abelmalu/golang-posts/Chat/internal/broker"
	"github.com/abelmalu/golang-posts/Chat/internal/handlers"
	"github.com/abelmalu/golang-posts/Chat/internal/kafka"
	miniocl "github.com/abelmalu/golang-posts/Chat/internal/minio"
	"github.com/abelmalu/golang-posts/Chat/internal/repository"
	"github.com/abelmalu/golang-posts/Chat/internal/service"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type App struct {
	mongoClient *mongo.Database
	router      *gin.Engine
	minioClient *minio.Client
}

var logger = platform.InitZapLogger()

func NewApp() *App {

	cfg, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("Error loading environment variables %v", err)
	}

	DB := initDB(cfg)

	r := gin.New()

	minio, err := miniocl.NewMinioClient()

	return &App{

		mongoClient: DB,
		router:      r,
		minioClient: minio,
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

func InitRoute(h *handlers.ChatHandler, r *gin.Engine) {

	r.Handle(http.MethodGet, "/api/chat/ws", h.HandleWebSocket)
	r.Handle(http.MethodGet, "/api/chat/user/chats", h.GetUserChats)

}

func (app *App) Run() {

	cb := broker.NewChatBroker(logger)
	cr := repository.NewChatRepository(app.mongoClient)
	cs := service.NewChatService(cr, app.minioClient)
	ch := handlers.NewChatHandler(cs, logger, cb)

	var broker = []string{"localhost:9092"}
	var userCreatedTopic = "userCreated"

	go kafka.StartConsumer(broker, userCreatedTopic, cs, logger)

	InitRoute(ch, app.router)

	log.Println("Starting server on port", "8084")
	app.router.Run(":8084")

}

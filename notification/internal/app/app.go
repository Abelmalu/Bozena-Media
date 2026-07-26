package application

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/abelmalu/golang-posts/notification/config"
	"github.com/abelmalu/golang-posts/notification/internal/broker"
	"github.com/abelmalu/golang-posts/notification/internal/handler"
	"github.com/abelmalu/golang-posts/notification/internal/kafka"
	"github.com/abelmalu/golang-posts/notification/internal/repository"
	"github.com/abelmalu/golang-posts/notification/internal/service"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)

type App struct {
	DB     *sql.DB
	Config *config.Config
}

var logger = platform.InitZapLogger()

func NewApp() *App {

	config, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("Couldn't load configuration %v", err)
	}

	DBConnPool, err := initDB(config)
	if err != nil {

		log.Fatalf("error while connectiong DB %v", err)
	}

	return &App{

		DB:     DBConnPool,
		Config: config,
	}
}

func initDB(config *config.Config) (*sql.DB, error) {

	DBUrl := config.DBURL

	DBConnPool, err := sql.Open("pgx", DBUrl)

	if err != nil {

		return nil, err
	}
	DBConnPool.SetMaxOpenConns(25)
	DBConnPool.SetMaxIdleConns(10)
	DBConnPool.SetConnMaxLifetime(5 * time.Minute)

	if err := DBConnPool.Ping(); err != nil {

		return nil, fmt.Errorf("pinging %s database: %w", "pgx", err)

	}
	logger.Info("Database connected successfully!")

	return DBConnPool, nil

}

func (app *App) Run() {

	notificationBroker := broker.NewNotificationBroker(logger)
	notificationRepo := repository.NewNotificationRepository(app.DB)
	notificationService := service.NewNotificationService(notificationRepo)
	notificationHandler := handler.NewNotificationHandler(logger,notificationService,notificationBroker)

	router := InitRoute(notificationHandler)

	brokers := []string{"localhost:9092"}
	followedtopic := "followed"
	userCreatedTopic := "userCreated"

	go kafka.StartConsumer(brokers, userCreatedTopic,followedtopic,notificationService,logger,notificationBroker)

	if err := router.Run(":8083"); err != nil {

		log.Fatalf("Error starting server %v", err)
	}

}

func InitRoute(handler *handler.NotificationHanlder) *gin.Engine {

	router := gin.New()
	//router.GET("/api/notification/stream",handler.Stream)
	router.Handle(http.MethodGet, "/api/notification/stream", handler.Stream)
	router.Handle(http.MethodGet, "/api/notification/user", handler.GetUserNotifications)


	return router

}

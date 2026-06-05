package application

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/abelmalu/golang-posts/like/config"
	handler "github.com/abelmalu/golang-posts/like/internal/handlers"
	"github.com/abelmalu/golang-posts/like/internal/kafka"
	"github.com/abelmalu/golang-posts/like/internal/repository"
	"github.com/abelmalu/golang-posts/like/internal/service"
	"github.com/abelmalu/golang-posts/like/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

type App struct {
	DB     *sql.DB
	config *config.Config
}

var logger = platform.InitZapLogger()

func NewApp() *App {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Couldn't load configuration %v", err)

	}

	DBConPool, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Error while initiating db connection %v", err)
	}

	return &App{
		DB:     DBConPool,
		config: cfg,
	}

}

func initDB(config *config.Config) (*sql.DB, error) {

	dsn := config.DBURL

	// creating the connection pool
	DBConPool, err := sql.Open("pgx", dsn)

	if err != nil {

		return nil, err
	}
	DBConPool.SetMaxOpenConns(25)
	DBConPool.SetMaxIdleConns(10)
	DBConPool.SetConnMaxLifetime(5 * time.Minute)

	// Check if connection credentials are correct
	if err := DBConPool.Ping(); err != nil {

		return nil, fmt.Errorf("pinging %s database: %w", "pgx", err)

	}
	logger.Info("Database connected successfully!")

	return DBConPool, nil
}

func (app *App) Run() {

	lis, _ := net.Listen("tcp", ":50053")
	s := grpc.NewServer()

	logger := platform.InitZapLogger()

	likeRepository := repository.NewLikeRepository(app.DB)
	likeService := service.NewLikeRepository(likeRepository)
	likeHandler := handler.NewLikeHandler(likeService, logger)

	pb.RegisterLikeServiceServer(s, likeHandler)

	brokers := []string{"localhost:9092"}
	Usertopic := "userCreated"
	postTopic := "postCreated"


	// Run consumer in the background
	go kafka.StartConsumer(brokers, Usertopic,postTopic, likeService, logger)

	s.Serve(lis)

}

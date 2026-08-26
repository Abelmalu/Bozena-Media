package application

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"
	"github.com/abelmalu/golang-posts/Feed/config"
	"github.com/abelmalu/golang-posts/Feed/internal/handler"
	"github.com/abelmalu/golang-posts/Feed/internal/kafka"
	miniocl "github.com/abelmalu/golang-posts/Feed/internal/minio"
	"github.com/abelmalu/golang-posts/Feed/internal/repository"
	"github.com/abelmalu/golang-posts/Feed/internal/service"
	"github.com/abelmalu/golang-posts/Feed/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	DB     *sql.DB
	config *config.Config
}

func NewApp() *App {

	cfg, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("error while loading env variables %v", err)
	}

	DBConn, err := initDB(cfg)

	if err != nil {

		log.Fatalf("Error while DB Connection")

	}

	return &App{
		DB:     DBConn,
		config: cfg,
	}
}

func initDB(cfg *config.Config) (*sql.DB, error) {

	dsn := cfg.DBURL

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

var logger = platform.InitZapLogger()

func (app *App) Run() {

	lis, _ := net.Listen("tcp", app.config.GRPCPORT)

	s := grpc.NewServer()

	minioClient, err := miniocl.NewMinioClient(app.config.MinioADD)

	if err != nil {

		logger.Error("minion Error", zap.Error(err))
	}
	feedRepo := repository.NewFeedRepository(app.DB)
	feedService := service.NewFeedService(feedRepo, minioClient)
	fh := handler.NewFeedHandler(feedService, logger)

	pb.RegisterFeedServiceServer(s, fh)

	Usertopic := "userCreated"
	postTopic := "postCreated"
	likedTopic := "liked"
	unLikedTopic := "unliked"
	followedTopic := "followed"
	unfollowedTopic := "unfollowed"

	// Run consumer in the background
	go kafka.StartConsumer(app.config.KafkaBrokersURL, Usertopic, postTopic, likedTopic, unLikedTopic, followedTopic, unfollowedTopic, feedService, logger)

	s.Serve(lis)

}

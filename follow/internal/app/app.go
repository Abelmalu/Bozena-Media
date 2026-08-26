package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/follow/config"
	handler "github.com/abelmalu/golang-posts/follow/internal/handlers"
	"github.com/abelmalu/golang-posts/follow/internal/kafka"
	"github.com/abelmalu/golang-posts/follow/internal/repository"
	"github.com/abelmalu/golang-posts/follow/internal/service"
	"github.com/abelmalu/golang-posts/follow/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	DB          *sql.DB
	Config      *config.Config
	RedisClient *redis.Client
	KafkaClient sarama.SyncProducer
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
	// Init Reids
	redisClient := initRedis(config.RedisADD, "", 0, logger)

	// initializing the kafka client

	kafkaClient, err := kafka.InitKafkaProducer(logger,config)

	if err != nil {

		log.Fatalf("Couldn't make kafka connection %v", err)
	}

	return &App{
		DB:          DBConnPool,
		Config:      config,
		RedisClient: redisClient,
		KafkaClient: kafkaClient,
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

// initRedis initializes redis client for the application
func initRedis(address, password string, db int, logger *platform.Logger) *redis.Client {
	ctx := context.Background()

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		},
	)

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	logger.Info("Redis connected successfully!", zap.String("Response", pong))

	return redisClient

}

func (app *App) Run() {

	port := app.Config.GRPCPORT
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {

		log.Fatalf("error while tcp connection %v", err)
	}
	s := grpc.NewServer()

	followRepo := repository.NewFollowRepository(app.DB)
	followService := service.NewFollowService(followRepo, app.KafkaClient)
	fh := handler.NewFollowHandler(followService, logger)

	pb.RegisterFollowServiceServer(s, fh)

	topic := "userCreated"

	// Run consumer in the background
	go kafka.StartConsumer(app.Config.KafkaBrokersURL, topic, followService, logger)

	s.Serve(lis)

}

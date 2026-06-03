package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/Auth/config"
	handler "github.com/abelmalu/golang-posts/Auth/internal/handlers"
	"github.com/abelmalu/golang-posts/Auth/internal/kafka"
	"github.com/abelmalu/golang-posts/Auth/internal/repository"
	"github.com/abelmalu/golang-posts/Auth/internal/service"
	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	config *config.Config
	DB     *sql.DB
	RedisClient *redis.Client
	Kafka   sarama.SyncProducer
}


	var logger = platform.InitZapLogger()

// NewApp creates the application instance  
func NewApp() (*App, error) {

	config, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("Couldn't load configuration %v", err)
	}

	DBConPool, err := initDB(config)
	if err != nil {

		log.Fatalf("Error while initiating db connection %v", err)

	}

	redisClient := initRedis("127.0.0.1:6379","",0,logger)


	// initializing the kafka client 

	kafkaClient,err := kafka.InitKafkaProducer(logger)

	if err != nil {

		log.Fatalf("Couldn't make kafka connection %v",err)
	}

	app := App{
		config: config,
		DB:     DBConPool,
		RedisClient: redisClient,
		Kafka: kafkaClient,
	}

	return &app, nil

}

// initDB initialize the database connection for the service
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

//initRedis initializes redis client for the application
func initRedis(address,password string,db int,logger *platform.Logger) (*redis.Client){
	ctx := context.Background()

	redisClient := redis.NewClient(
		&redis.Options{
			Addr: address,
			Password: password,
			DB: db,
		},
	)

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	logger.Info("Redis connected successfully!",zap.String("Response",pong))

	return redisClient


}

// Run starts the gRPC server on the provided port
func (app *App) Run() {

	lis, _ := net.Listen("tcp", ":50052")
	s := grpc.NewServer()

	
    // Dependency Injection for each layer one by one 
	authRepo := repository.NewAuthRepository(app.DB)
	authService := service.NewAuthService(authRepo,app.RedisClient,app.Kafka,logger)
	authHandler := handler.NewAuthHandler(authService,logger)

	
	pb.RegisterAuthServiceServer(s, authHandler)
	s.Serve(lis)
	

}

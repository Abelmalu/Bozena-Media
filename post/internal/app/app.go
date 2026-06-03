package application

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/abelmalu/golang-posts/platform"
	"github.com/abelmalu/golang-posts/post/config"
	"github.com/abelmalu/golang-posts/post/internal/handlers"
	"github.com/abelmalu/golang-posts/post/internal/interceptors"
	"github.com/abelmalu/golang-posts/post/internal/kafka"
	"github.com/abelmalu/golang-posts/post/internal/repository"
	"github.com/abelmalu/golang-posts/post/internal/service"
	"github.com/abelmalu/golang-posts/post/proto/pb"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

type App struct {
	config *config.Config
	DB     *sql.DB
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

	app := App{
		config: config,
		DB:     DBConPool,
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

func (app *App) Run() {
	logger := platform.InitZapLogger()


	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        interceptors.AuthInterceptor(logger),         
        interceptors.PostOwnershipInterceptor(app.DB,logger), // checks if the user is post owner
    ),
)


	
	postRepo := repository.NewPostRepository(app.DB)
	postService := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService,logger)

	
	pb.RegisterPostServiceServer(s, postHandler)


	brokers := []string{"localhost:9092"}
	topic := "test-topic"

	// Run consumer in the background
	go kafka.StartConsumer(brokers, topic,postService,logger)

	// start the grpc server
	s.Serve(lis)
	

}

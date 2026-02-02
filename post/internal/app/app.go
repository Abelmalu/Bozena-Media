package application

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/abelmalu/golang-posts/post/internal/handlers"
	"github.com/abelmalu/golang-posts/post/config"
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

type postServer struct {
	pb.UnimplementedPostServiceServer
}

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

	return DBConPool, nil
}

// Run starts the gRPC server on the provided port
func (app *App) Run() {

	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	
    // Dependency Injection for each layer one by one 
	postRepo := repository.NewPostRepository(app.DB)
	postService := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	
	pb.RegisterPostServiceServer(s, postHandler)
	s.Serve(lis)
	

}

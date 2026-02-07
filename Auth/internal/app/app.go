package application

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/abelmalu/golang-posts/Auth/config"
	handler "github.com/abelmalu/golang-posts/Auth/internal/handlers"
	"github.com/abelmalu/golang-posts/Auth/internal/repository"
	"github.com/abelmalu/golang-posts/Auth/internal/service"
	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

type App struct {
	config *config.Config
	DB     *sql.DB
}

type postServer struct {
	pb.UnimplementedAuthServiceServer
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

	lis, _ := net.Listen("tcp", ":50052")
	s := grpc.NewServer()
	
    // Dependency Injection for each layer one by one 
	authRepo := repository.NewAuthRepository(app.DB)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	
	pb.RegisterAuthServiceServer(s, authHandler)
	s.Serve(lis)
	

}

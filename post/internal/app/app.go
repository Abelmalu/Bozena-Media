package application

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/abelmalu/golang-posts/post/config"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

type App struct {
	config     *config.Config
	grpcServer *grpc.Server
	DB         *sql.DB
}

// NewApp creates the application instance
func NewApp() (*App, error) {

	config, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("Couldn't load configuration %v", err)
	}

	DBConPool, err = initDB(config)

	return nil, nil

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
func (app *App) Run(){

		grpcServer := grpc.NewServer()




}

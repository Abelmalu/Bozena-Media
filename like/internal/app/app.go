package application

import (
	"database/sql"
	"log"
	"fmt"
	"time"
	"github.com/abelmalu/golang-posts/like/config"
)

type App struct {
	DB *sql.DB
	config *config.Config
}

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

	return DBConPool, nil
}


func (app *App) Run() {

	

}


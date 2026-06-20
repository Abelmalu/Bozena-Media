package application

import (
	"database/sql"
	"fmt"
	"log"
//	"net"
	"time"

	"github.com/abelmalu/golang-posts/Feed/config"
	"github.com/abelmalu/golang-posts/platform"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	//"google.golang.org/grpc"
)

type App struct {
	DB     *sql.DB
	config *config.Config
}

var logger = platform.InitZapLogger()

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




// func (app *App) Run () {

// 	lis,_ := net.Listen("tcp",":50055")

// 	s := grpc.NewServer()
	
// }



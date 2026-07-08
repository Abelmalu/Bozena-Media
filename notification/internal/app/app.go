package application

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/abelmalu/golang-posts/notification/config"
	"github.com/abelmalu/golang-posts/platform"
)



type App struct {

	DB *sql.DB
	Config  *config.Config
 
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


	return &App{

		DB: DBConnPool,
		Config: config,
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

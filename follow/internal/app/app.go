package application

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/abelmalu/golang-posts/follow/config"
	"github.com/abelmalu/golang-posts/follow/internal/handlers"
	"github.com/abelmalu/golang-posts/follow/internal/repository"
	"github.com/abelmalu/golang-posts/follow/internal/service"
	"github.com/abelmalu/golang-posts/follow/proto/pb"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type App struct {
	DB          *sql.DB
	Config      *config.Config
	RedisClient *redis.Client
}

func NewApp() *App {

	config, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("Couldn't load configuration %v", err)
	}

	DBConnPool, err := initDB(config)

	return &App{
		DB: DBConnPool,
	}
}

func initDB(config *config.Config) (*sql.DB, error) {

	DBUrl := config.DBURL

	DBConnPool, err := sql.Open(DBUrl, "pgx")

	if err != nil {

		return nil, err
	}
	DBConnPool.SetMaxOpenConns(25)
	DBConnPool.SetMaxIdleConns(10)
	DBConnPool.SetConnMaxLifetime(5 * time.Minute)

	if err := DBConnPool.Ping(); err != nil {

		return nil, fmt.Errorf("pinging %s database: %w", "pgx", err)

	}

	return DBConnPool, nil

}

func (app *App) Run() {
	port := app.Config.GRPCPORT
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%v",port))
	s := grpc.NewServer()

	followRepo := repository.NewFollowRepository(app.DB)
	followService := service.NewFollowService(followRepo)
	followHandler := handler.NewFollowHandler(followService)

	pb.RegisterFollowServiceServer(s,followHandler)
	s.Serve(lis)


}

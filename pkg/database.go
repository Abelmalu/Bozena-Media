package pkg

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)



//function InitDB establishes a connection pool for the entire application
func InitDB()(*sql.DB,error){

	// connection string for postgresql 
	dsn := "postgres://root:passworD-123@localhost:5432/blog?sslmode=disable"

     // creating the connection pool
	dbConPool,err :=sql.Open("pgx",dsn)

	if err != nil{

		return nil,err
	}
	dbConPool.SetMaxOpenConns(25)
	dbConPool.SetMaxIdleConns(10)
	dbConPool.SetConnMaxLifetime(5*time.Minute)


	// Check if connection credentials are correct 
	if err := dbConPool.Ping(); err != nil{
		dbConPool.Close()

		return nil, fmt.Errorf("pinging %s database: %w", "pgx", err)

	}


   return dbConPool,nil
}
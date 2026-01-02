package pkg

import "database/sql"

var db *sql.DB


func initDB(){

	dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"


	dbConPool,err :=sql.Open("pgx",dsn)
}
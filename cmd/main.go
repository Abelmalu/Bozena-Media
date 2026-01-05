package main

import (
	"log"

	"github.com/abelmalu/golang-posts/internal"
	"github.com/abelmalu/golang-posts/pkg"
)


func main(){

	db,err := pkg.InitDB()

	if err != nil{

		log.Fatalf("Service startup failed: %v",err)
	}

	defer db.Close()

	router := internal.SetupRoutes()
	if err :=  router.Run(":8080"); err !=nil{

		log.Println("The router error is ",err)
	}
  

}

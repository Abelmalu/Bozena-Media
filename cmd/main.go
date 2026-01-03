package main

import (
	"log"

	"github.com/abelmalu/golang-posts/pkg"
)


func main(){

	db,err := pkg.InitDB()

	if err != nil{

		log.Fatalf("Service startup failed: %v",err)
	}

	defer db.Close()
  

}

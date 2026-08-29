package server

import (
	"log"
	"github.com/gin-gonic/gin"
)


func StartServer(router *gin.Engine,port string){

	
	if err := router.Run(port); err != nil {

		log.Fatalf("Couldn't start the router %v",err)
	}


}
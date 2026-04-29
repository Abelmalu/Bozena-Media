package server

import (
	"log"

	"github.com/gin-gonic/gin"
)


func StartServer(router *gin.Engine){

	
	if err := router.Run(":8082"); err != nil {

		log.Fatalf("Couldn't star the router %v",err)
	}


}
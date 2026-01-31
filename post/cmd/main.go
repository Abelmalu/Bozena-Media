package main

import (
	"context"
	"log"
	"net"

	application "github.com/abelmalu/golang-posts/post/internal/app"
	"github.com/abelmalu/golang-posts/post/proto/pb"
	"google.golang.org/grpc"
)

func main() {
    //Initiating the application

    app,err := application.NewApp()
    if err != nil{

        log.Fatalf("Application Initializing Error %v",err)
    }

    app.Run()

   
}



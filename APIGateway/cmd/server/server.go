package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {

	return gin.New()

}
func StartServer(router *gin.Engine, port string) {

	srv := &http.Server{

		Addr:    port,
		Handler: router,
	}

	serverErrors := make(chan error, 1)

	shutDownSignal := make(chan os.Signal, 1)

	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {

			serverErrors <- err
		}
	}()

	select {

	case err := <-serverErrors:

		log.Fatalf("Server error happend: %v", err)

	case  <-shutDownSignal:

		shutDownCtx,cancel := context.WithTimeout(context.Background(),time.Second * 5)
		defer cancel()

		if err := srv.Shutdown(shutDownCtx); err != nil {


			log.Fatalf("Server forced to shutdwon %v",err)
		}


		log.Println("Server exited cleanly.")



	}

}

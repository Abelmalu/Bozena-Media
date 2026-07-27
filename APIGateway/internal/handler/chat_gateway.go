package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	logger *platform.Logger
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewChatHandler(l *platform.Logger) *ChatHandler {

	return &ChatHandler{

		logger: l,
	}

}

func (ch *ChatHandler) Connect(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}
	defer conn.Close()

	log.Println("Client successfully connected!")

	for {
		messageType, messagePayload, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Connection closed or read error: %v\n", err)
			break
		}

		fmt.Printf("Received: %s\n", messagePayload)

		err = conn.WriteMessage(messageType, messagePayload)
		if err != nil {
			log.Printf("Failed to send message: %v\n", err)
			break
		}
	}

}

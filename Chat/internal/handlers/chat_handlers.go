package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abelmalu/golang-posts/Chat/internal/core"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	logger   *platform.Logger
	cs       core.ChatService
	WSUpGrader *websocket.Upgrader
}

func NewChatHandler(cs core.ChatService, lg *platform.Logger) *ChatHandler {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &ChatHandler{

		logger: lg,
		cs:     cs,
		WSUpGrader:&upgrader,
	}
}

func (ch *ChatHandler) HandleWebSocket(c *gin.Context) {
	conn, err := ch.WSUpGrader.Upgrade(c.Writer, c.Request, nil)
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

func (ch *ChatHandler) CreateConversation(){}

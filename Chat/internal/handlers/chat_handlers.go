package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	//ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/Chat/internal/broker"
	"github.com/abelmalu/golang-posts/Chat/internal/core"
//	"github.com/abelmalu/golang-posts/notification/pkg"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ChatHandler struct {
	logger     *platform.Logger
	cs         core.ChatService
	WSUpGrader *websocket.Upgrader
	broker     *broker.ChatBroker
}

func NewChatHandler(cs core.ChatService, lg *platform.Logger, b *broker.ChatBroker) *ChatHandler {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &ChatHandler{

		logger:     lg,
		cs:         cs,
		WSUpGrader: &upgrader,
		broker:     b,
	}
}

func (ch *ChatHandler) HandleWebSocket(c *gin.Context) {

    userID := c.GetHeader("X-User-ID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		ch.logger.Error("Error parsing user ID", zap.Error(err))
		//pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)
		return
	}
	log.Println("User connected:", userIDInt)

	conn, err := ch.WSUpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}
	defer conn.Close()
	log.Println("Client successfully connected!")

	userChan := ch.broker.Register(userIDInt)
	defer ch.broker.Unregister(userIDInt)

	readerDone := make(chan struct{})

	go func() {
		defer close(readerDone)
		for {
			_, messagePayload, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Connection closed or read error: %v\n", err)
				return
			}
			fmt.Printf("Received: %s\n", messagePayload)

			// Echo message back to client
			// err = conn.WriteMessage(messageType, messagePayload)
			// if err != nil {
			// 	log.Printf("Failed to send message: %v\n", err)
			// 	return
			// }

			msgStr := string(messagePayload)
			receverID, err := strconv.Atoi(msgStr)
	if err != nil {
		ch.logger.Error("Error parsing user ID", zap.Error(err))
		//pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)
		return
	}

			go func(msg string) {
				ch.broker.Publish((receverID), msg)
				fmt.Printf("after: %s\n", msg)
			}(string(messagePayload))

		}
	}()

	for {
		select {
		case <-c.Request.Context().Done():
			log.Println("Request context cancelled")
			return

		case <-readerDone:
			log.Println("Reader goroutine exited, closing handler")
			return

		case msg, ok := <-userChan:
			if !ok {
				log.Println("Broker channel closed")
				return
			}
			// Send broker messages as WebSocket Text messages instead of SSE
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("Failed to push broker message to client: %v\n", err)
				return
			}

		

		}
	}
}

func (ch *ChatHandler) CreateConversation() {

	

}

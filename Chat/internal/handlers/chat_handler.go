package handlers

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	chatclient "github.com/abelmalu/golang-posts/Chat/ChatClient"
	"github.com/abelmalu/golang-posts/Chat/internal/broker"
	"github.com/abelmalu/golang-posts/Chat/internal/core"
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

	requestID := c.GetHeader("X-Request-ID")
	ch.logger.Info(requestID)
	userID := c.GetHeader("X-User-ID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		ch.logger.Error("Error parsing user ID", zap.Error(err))
		//pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)
		return
	}

	conn, err := ch.WSUpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}
	log.Println("Client successfully connected!")

	client := &chatclient.Client{
		UserID:      userIDInt,
		Conn:        conn,
		Broker:      ch.broker,
		ChatService: ch.cs,
		Logger:      ch.logger,
	}

	userChan := client.Broker.Register(client.UserID)

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {

		defer wg.Done()

		client.ReadPump(c.Request.Context())

	}()

	go func() {

		defer wg.Done()

		client.WritePump(c.Request.Context(), userChan)

	}()

	wg.Wait()

}


func (ch *ChatHandler) GetUserChats(c *gin.Context) {

	requestID := c.GetHeader("X-Request-ID")
	userID := c.GetHeader("X-User-ID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		ch.logger.Error("Error parsing user ID", zap.Error(err))
		return
	}

	
	




	
}

package chatclient

import (
	"context"
	"fmt"
	"log"
	//"time"

	"github.com/abelmalu/golang-posts/Chat/internal/broker"
	"github.com/abelmalu/golang-posts/Chat/internal/core"
	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID int

	Conn *websocket.Conn

	Broker *broker.ChatBroker

	ChatService core.ChatService

	Logger *platform.Logger
}

func NewChatClient(userID int, conn *websocket.Conn, broker *broker.ChatBroker, ch core.ChatService, l *platform.Logger) *Client {

	return &Client{

		UserID:      userID,
		Conn:        conn,
		Broker:      broker,
		ChatService: ch,
		Logger:      l,
	}
}
func (c *Client) ReadPump(ctx context.Context) {

	defer func() {
		c.Broker.Unregister(c.UserID)
		c.Conn.Close()

	}()

	for {
		//ctx, cancel := context.WithTimeout(ctx, time.Second*1000)

		var req dto.MessageRequest

		if err := c.Conn.ReadJSON(&req); err != nil {

		//	cancel()

			return
		}

		// if err := c.ChatService.SendMessage(ctx, c.UserID, req.ReceiverID, req.Message); err != nil {

		// 	c.Logger.Error("Error saving chat message to db")
		// 	cancel()

		// 	return

		// }

		fmt.Println(req.ReceiverID)

		c.Broker.Publish(req.ReceiverID, req.Message)


	}
}

func (c *Client) WritePump(ctx context.Context, userChan chan string) {



	for {

		select {

		case msg, ok := <-userChan:

			if !ok {
				c.Logger.Error("error reading from channel")
				return
			}

			err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("Failed to push broker message to client: %v\n", err)
				return
			}

		case <-ctx.Done():


			return
		}
	}
}

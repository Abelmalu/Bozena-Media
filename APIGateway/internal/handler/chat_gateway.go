package handler

import (
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

func (h *ChatHandler) Connect(c *gin.Context) {

	clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer clientConn.Close()

	backendConn, _, err := websocket.DefaultDialer.Dial(
		"ws://localhost:8084/api/chat/ws",
		http.Header{
			"X-User-ID": []string{c.GetHeader("X-User-ID")},
		},
	)
	if err != nil {
		return
	}
	defer backendConn.Close()

	// Browser -> Chat Service
	go func() {
		for {
			mt, msg, err := clientConn.ReadMessage()
			if err != nil {
				backendConn.Close()
				return
			}

			if err := backendConn.WriteMessage(mt, msg); err != nil {
				return
			}
		}
	}()

	// Chat Service -> Browser
	for {
		mt, msg, err := backendConn.ReadMessage()
		if err != nil {
			break
		}

		if err := clientConn.WriteMessage(mt, msg); err != nil {
			break
		}
	}
}

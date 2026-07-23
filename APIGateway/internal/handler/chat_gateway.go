package handler

import (
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
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


	
}
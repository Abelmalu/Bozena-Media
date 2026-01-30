package handlers

import (
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/clients"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postClient *clients.PostClient
}

func NewPostHandler(pc *clients.PostClient) *PostHandler {
	return &PostHandler{postClient: pc}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req struct {
		UserID  int64 `json:"user_id"`
		Title string `json:"title"`
		Content string `json:"content"`
	}

	if err :=c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"status":"error",
			"message":"Bad Request",

		})

		return
	}

	resp, err := h.postClient.CreatePost(req.UserID, req.Title,req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":"error",
			"message":"Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"status":"success",
		"message":resp,
	})
}
package handlers

import (
	"log"
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/clients"
	"github.com/abelmalu/golang-posts/internal/models"
	"github.com/abelmalu/golang-posts/pkg"
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
		UserID  int64  `json:"user_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
		})

		return
	}

	resp, err := h.postClient.CreatePost(req.UserID, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

//ListPosts get all posts from the db

func GetPosts(c *gin.Context) {
	var posts []models.Post
	query := `SELECT * FROM posts`

	rows, err := pkg.DB.Query(query)
	if err != nil {

		log.Printf("error during db query %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post

		rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserID)
		posts = append(posts, post)

	}

	if err = rows.Err(); err != nil {
		log.Printf("error after iterating rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, posts)
}

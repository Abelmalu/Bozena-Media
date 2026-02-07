package handler

import (
	"context"
	"net/http"
	"github.com/abelmalu/golang-posts/post/proto/pb"
	"github.com/gin-gonic/gin"
)

type PostService interface{
	CreatePost(ctx context.Context,userID int64, title, content string) (*pb.CreatePostResponse, error)
	ListPosts()(*pb.ListPostsResponse,error)
	UpdatePost (postID int64, title, content string)(*pb.UpdatePostResponse,error)
	DeletePost (postID int64)(*pb.DeletePostResponse,error)
}

type PostHandler struct {
	postClient PostService
}

func NewPostHandler(pc PostService) *PostHandler {
	return &PostHandler{postClient: pc}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	
	var req struct {
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
	userIDValue,exists := c.Get("userID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return

	}
	userID := userIDValue.(int64)

	

	resp, err := h.postClient.CreatePost(c.Request.Context(),userID, req.Title, req.Content)
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

// func ListPosts(c *gin.Context) {
	
// 	query := `SELECT * FROM posts`

// 	rows, err := pkg.DB.Query(query)
// 	if err != nil {

// 		log.Printf("error during db query %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"message": "Internal Server Error",
// 		})
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var post models.Post

// 		rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserID)
// 		posts = append(posts, post)

// 	}

// 	if err = rows.Err(); err != nil {
// 		log.Printf("error after iterating rows: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"message": "Internal Server Error",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, posts)
// }

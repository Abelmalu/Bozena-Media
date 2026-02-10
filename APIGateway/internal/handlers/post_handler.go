package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/abelmalu/golang-posts/post/proto/pb"
	"github.com/gin-gonic/gin"
)

type PostService interface{
	CreatePost(ctx context.Context,userID int64, title, content string) (*pb.CreatePostResponse, error)
	ListPosts()(*pb.ListPostsResponse,error)
	UpdatePost (ctx context.Context,postID int64,title string,content string)(*pb.UpdatePostResponse,error)
	DeletePost (ctx context.Context,postID int64)(*pb.DeletePostResponse,error)
}

type PostHandler struct {
	postClient PostService
}

func NewPostHandler(pc PostService) *PostHandler {
	return &PostHandler{postClient: pc}
}

func (postHandler *PostHandler) CreatePost(c *gin.Context) {
	
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
	userIDInt := userIDValue.(int)
	userID := int64(userIDInt)

	resp, err := postHandler.postClient.CreatePost(c.Request.Context(),userID, req.Title, req.Content)
	if err != nil {

		log.Printf("the error is %v",err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (postHandler *PostHandler) UpdatePost(c *gin.Context){

	var input struct {
		Title   string `json:"title" db:"title"`
		Content string `json:"content" db:"content"`
	}
	postIDStr := c.Param("id")
	postIDValue, err := strconv.Atoi(postIDStr)
	postID := int64(postIDValue)

	if err != nil {

		log.Printf("error while Atoi, %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return

	}
	if err := c.ShouldBindJSON(&input); err != nil {

		log.Printf("error while parsing JSON %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
		})
		return
	}

	resp,err := postHandler.postClient.UpdatePost(c.Request.Context(),postID,input.Title,input.Content)
	if err != nil{


		log.Printf("error from post service %v",err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":"error",
			"message":"Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusCreated,resp)
}
func (postHandler *PostHandler) DeletePost(c *gin.Context){

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)

	if err != nil {

		log.Printf("error while Atoi, %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
		})
		return

	}

	resp, err := postHandler.postClient.DeletePost(c.Request.Context(),int64(postID))

	if err != nil{

		log.Printf("error from post service %v",err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":"error",
			"message":"Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusAccepted, resp)
}

func (postHandler *PostHandler) ListPosts(c *gin.Context) {

	resp,err := postHandler.postClient.ListPosts()

	if err != nil{

		log.Printf("the error is %v",err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":"error",
			"message":"Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK,resp)



}

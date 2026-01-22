package posts

import (
	"log"
	"net/http"
	"strconv"

	"github.com/abelmalu/golang-posts/internal/models"
	"github.com/abelmalu/golang-posts/pkg"
	"github.com/gin-gonic/gin"
)

func Posts(c *gin.Context) {

}

func CreatePost(c *gin.Context) {
	//Grabbing the value of userID from the context

	log.Printf("in the handler")
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}

	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {

		log.Printf("JSON binding error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}
	post.UserID = userIDValue.(int)
	query := `INSERT INTO posts (title,content,user_id) VALUES($1,$2,$3)`

	_, err := pkg.DB.Exec(query, post.Title, post.Content, post.UserID)

	if err != nil {
		log.Printf("database exec error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully created a post",
	})
	return

}

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

func UpdatePost(c *gin.Context) {

	var input struct {
		Title   string `json:"title" db:"title"`
		Content string `json:"content" db:"content"`
	}
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)

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

	query := `UPDATE posts SET title=$1, content=$2 WHERE id=$3`

	result, err := pkg.DB.Exec(query, input.Title, input.Content, postID)
	if err != nil {
		log.Printf("DB ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("DB error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ooops Couldn't update post!",
		})
		return

	}
	if rowsAffected == 0 {
		log.Printf("DB erro rows affected came zero")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ooops Couldn't update post!",
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Post updated Successfully",
	})

}

func DeletePost(c *gin.Context){

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

	query := `DELETE FROM posts WHERE id=$1`

	result,err := pkg.DB.Exec(query,postID)
	
	if err != nil{
		log.Printf("DB erro %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return

	}
	rowsAffected,err := result.RowsAffected()
	if err != nil{

		log.Printf("DB erro %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ooops Couldn't delete post!",
		})
		return

	}

	if rowsAffected == 0{
		log.Printf("DB erro rows affected came zero")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ooops Couldn't delete post!",
		})
		return


	}
	c.JSON(http.StatusOK,gin.H{
		"status":"success",
		"message":"Post deleted successfully",
	})


}

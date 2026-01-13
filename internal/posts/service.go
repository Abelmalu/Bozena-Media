package posts

import (
	"log"
	"net/http"
	"github.com/abelmalu/golang-posts/internal/models"
	"github.com/abelmalu/golang-posts/pkg"
	"github.com/gin-gonic/gin"
)




func CreatePost(c *gin.Context){
	//Grabbing the value of userID from the context 

	log.Printf("in the handler")
	 userIDValue, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return

    }

  

	var post models.Post

	if err:=c.ShouldBindJSON(&post); err != nil{

		log.Printf("JSON binding error %v",err)
		c.JSON(http.StatusBadRequest,gin.H{
			"status":"error",
			"message":"Invalid request body",
		
		})
		return
	}
     post.UserID = userIDValue.(int)
	query := `INSERT INTO posts (title,content,user_id) VALUES($1,$2,$3)`


	_,err := pkg.DB.Exec(query,post.Title,post.Content,post.UserID)

	if err != nil{
		log.Printf("database exec error %v",err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":"error",
			"message":"Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"status":"success",
		"message":"Successfully created a post",
	})
     return 

}
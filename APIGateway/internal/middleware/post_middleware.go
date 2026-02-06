package middleware

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)



func AuthorizePostOwner(DB *sql.DB) gin.HandlerFunc{

	return func(c *gin.Context) {
		userIDValue,exists := c.Get("userID")
		userID := userIDValue.(int)
		posIDStr := c.Param("id")
		postID,err := strconv.Atoi(posIDStr)
		if err != nil{
			log.Printf("error during string to int conversion %v",err)
			c.JSON(http.StatusBadRequest,gin.H{
				"status":"error",
				"message":"Bad Request",
			})
			return
		}
		

		if !exists{

			log.Printf("can't find userID from the request")
			c.JSON(http.StatusBadRequest,gin.H{
				"status":"error",
				"message":"Bad Request",
			})
		}

		var postOwnerID int 

		query := 	`SELECT user_id FROM posts WHERE id=$1`

		err = DB.QueryRow(query, postID).Scan(&postOwnerID)

      
        if err != nil {
            log.Printf("Post ownership check failed: %v", err)
            c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
            c.Abort()
            return
        }

		if userID != postOwnerID{

			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to modify this post"})
            c.Abort() // Stop the request here!
            return
		}

		c.Next()


		
	}


}
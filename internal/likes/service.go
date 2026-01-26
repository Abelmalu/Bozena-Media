package likes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateLike(c *gin.Context) {

	// Getting post id from the reques parameter
	postIDStr := c.Param("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {

		log.Printf("error while converting post id to integer %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Bad Request",
		})
		return

	}

	userIDValue,exists := c.Get("userID")

	if !exists {

		log.Printf("Can't get user id from the request")
		c.JSON(http.StatusBadRequest,gin.H{
			"status":"error",
			"message":"Bad Request",
		})
		return 
	}

	//asserting interface userIDvalue to integer
	userID := userIDValue.(int)


	// checking if the like exists then toggle 

	query := `SELECT * FROM likes WHERE post_id=$1, user_id=$2`


}

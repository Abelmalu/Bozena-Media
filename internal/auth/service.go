package auth

import (
	"net/http"
	"github.com/abelmalu/golang-posts/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/abelmalu/golang-posts/pkg"

)



func Register(c *gin.Context){

	var newUser models.User

	


	if err := c.ShouldBindJSON(&newUser); err !=nil{

		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return

	}

	  db,err := pkg.InitDB()

	  if err != nil{

		c.JSON(http.StatusInternalServerError,gin.H{"status":"error", "message":err.Error()})
		return
	  }

	  query := `INSERT INTO users(name,username,email,password) VALUES($1,$2,$3,$4)`
	  _,dbError := db.Exec(query,newUser.Name,newUser.Username,newUser.Email, newUser.Password)
	  if dbError != nil{

		c.JSON(http.StatusInternalServerError,gin.H{"status":"error", "message":dbError.Error()})
		return
	  }

	responseMessage := "Welcome " + newUser.Name
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": responseMessage,
		})


}


// login authenticate a user with a valid username and password
func Login( c *gin.Context){
  
	c.SholudBindJSON()




}
package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/abelmalu/golang-posts/internal/models"
	"github.com/abelmalu/golang-posts/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func Register(c *gin.Context) {

	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	db := pkg.DB

	query := `INSERT INTO users(name,username,email,password) VALUES($1,$2,$3,$4)`
	_, dbError := db.Exec(query, newUser.Name, newUser.Username, newUser.Email, newUser.Password)
	if dbError != nil {

		// Change *pq.Error to *pgconn.PgError
		if pgErr, ok := dbError.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"message": "Username or Email already exists"})
				return
			}
		}
		log.Printf("Registration DB Error %v", dbError)

		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal Server Error"})
		return
	}

	responseMessage := "Welcome " + newUser.Name
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": responseMessage,
	})

}

// login authenticate a user with a valid username and password
func Login(c *gin.Context) {

	var user models.User

	var input struct {
		Username string `json:"username" db:"username"`
		Password string `json:"password" db:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	query := `SELECT * FROM users where username=$1`
	err := pkg.DB.QueryRow(query,input.Username).Scan(&user.Id,&user.Name,&user.Username,&user.Password,&user.Email,&user.CreatedAt,&user.UpdatedAt)
	if err != nil{

		log.Printf("Login DB Error %v",err)
		c.JSON(http.StatusBadRequest,gin.H{"status":"error", "message":"Invalid Credentials"})
		return 
	}

	
	trimmedPassword := strings.TrimSpace(input.Password)

	if trimmedPassword == user.Password{

		log.Printf("user %s LoggedIn",user.Username)


	}else{

		c.JSON(http.StatusBadRequest,gin.H{"status":"error","message":"Invalid Credentials"})


	}

}

func Home(c *gin.Context){


	c.JSON(http.StatusOK,"welcome back")
}

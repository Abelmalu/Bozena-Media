package auth

import (
	"log"
	"net/http"

	"github.com/abelmalu/golang-posts/internal/models"
	"github.com/abelmalu/golang-posts/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {

	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newUser.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("Password hash error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}
	newUser.Password = string(hashedPassword)

	db := pkg.DB

	query := `INSERT INTO users(name,username,email,password) VALUES($1,$2,$3,$4) RETURNING id`
	if err := db.QueryRow(query, newUser.Name, newUser.Username, newUser.Email, newUser.Password).Scan(&newUser.ID); err != nil {
		// Change *pq.Error to *pgconn.PgError
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"message": "Username or Email already exists"})
				return
			}
		}
		log.Printf("Registration DB Error %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal Server Error"})
		return

	}

	token, err := pkg.GenerateAcessToken(newUser.ID)
	if err != nil {

		log.Fatalf("JWT error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server error",
		})
		return
	}

	responseMessage := "Welcome " + newUser.Name
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": responseMessage,
		"token":   token,
	})

}

// login authenticate a user with a valid username and password
func Login(c *gin.Context) {

	var user models.User

	 input := struct {
		Username string `json:"username" db:"username"`
		Password string `json:"password" db:"password"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	query := `SELECT * FROM users where username=$1`
	err := pkg.DB.QueryRow(query, input.Username).Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {

		log.Printf("Login DB Error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),  // stored hash
		[]byte(input.Password), // plain password
	)

	if err != nil {
		// password mismatch
		log.Printf("Incorrect password error %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}

	token, err := pkg.GenerateAcessToken(user.ID)
	if err != nil {

		log.Fatalf("JWT error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Successfully LoggedIn!",
		"token":   token,
	})

}

func Home(c *gin.Context) {

	c.JSON(http.StatusOK, "welcome back")
}

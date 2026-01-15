package auth

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/abelmalu/golang-posts/internal/models"
	"github.com/abelmalu/golang-posts/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

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

	// if registration is successful check request header and generate tokens
	clientHeader := c.GetHeader("X-Client-Type")
	clientType := models.ClientWeb

	if clientHeader == "mobile" {

		clientType = models.ClientMobile
	}

	tokens, err := issueTokens(c, newUser.ID, clientType)

	if err != nil {
		log.Fatalf("JWT error %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})

	}

	response := gin.H{"message": "Login successful"}

	switch clientType {
	case models.ClientWeb:
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			MaxAge:   30 * 24 * 60 * 60,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		c.Header("Authorization", "Bearer "+tokens.AccessToken)

	case models.ClientMobile:
		response["access_token"] = tokens.AccessToken
		response["refresh_token"] = tokens.RefreshToken
	}

	c.JSON(http.StatusOK, response)

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
	// if login is successful check request header and generate tokens
	clientHeader := c.GetHeader("X-Client-Type")
	clientType := models.ClientWeb

	if clientHeader == "mobile" {

		clientType = models.ClientMobile
	}

		tokens, err := issueTokens(c, user.ID, clientType)

	if err != nil {
		log.Fatalf("JWT error %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})

	}

	response := gin.H{"message": "Login successful"}

	switch clientType {
	case models.ClientWeb:
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			MaxAge:   30 * 24 * 60 * 60,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		c.Header("Authorization", "Bearer "+tokens.AccessToken)

	case models.ClientMobile:
		response["access_token"] = tokens.AccessToken
		response["refresh_token"] = tokens.RefreshToken
	}

	c.JSON(http.StatusOK, response)

}

func issueTokens(c *gin.Context, userID int, clientType models.ClientType) (*TokenPair, error) {

	accessToken, err := pkg.GenerateAcessToken(userID)
	if err != nil {

		return nil, err
	}
	refreshToken, err := pkg.GenerateAcessToken(userID)
	if err != nil {

		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func StoreRefreshTokens(userID int,refreshToken string, clientType string)(sql.Result,error){


	query := `INSERT INTO refresh_tokens (user_id,token_text,expires_at,client_type) VALUES($1,$2,$3,$4)`

	result,err := pkg.DB.Exec(query, userID,refreshToken,clientType)
	if err != nil{

		return nil,err
	}

	return result,nil

}
func ExtractRefreshToken(c *gin.Context) (string, error) {

	//Extracting  refresh tokens from mobile app clients 
	var  refreshToken string
	if err := c.ShouldBindJSON(&refreshToken); err == nil {
		if refreshToken != "" {
			return refreshToken, nil
		}
	}

	// Extracting HttpOnly cookie for web clients
	if token, err := c.Cookie("refresh_token"); err == nil && token != "" {
		return token, nil
	}

	//Extracting from custom header fallback
	if token := c.GetHeader("X-Refresh-Token"); token != "" {
		return token, nil
	}

	return "", errors.New("Refresh Token not found")
}
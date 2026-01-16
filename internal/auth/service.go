package auth

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

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

	// Generate tokens
	tokens, err := issueTokens(c, user.ID, clientType)

	if err != nil {
		log.Printf("JWT error %v", err)

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

// RefreshHandler handle requests for getting new access tokens
func RefreshHandler(c *gin.Context) {

	// extracting the refresh token from the request for both mobile and web clients
	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {

		log.Printf("refresh token extracting error %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return
	}
	// validate the token to check if it tampered
	_, err = pkg.ValidateRefreshToken(refreshToken)

	if err != nil {

		log.Printf("refresh token validation error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return
	}

	// Get the refresh token from the DB
	tokenRecord, err := GetRefreshToken(refreshToken)
	if err != nil {

		log.Printf("Getting refresh token database error %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return
	}
	// check if it is revoked or has expired
	if tokenRecord.Revoked || tokenRecord.ExpiresAt.Before(time.Now()) {
		log.Printf("token expired or revoked ")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return

	}
	var clientType models.ClientType
	userID := tokenRecord.UserID
	clientTypeStr := tokenRecord.ClientType

	switch clientTypeStr {
	case "web":
		clientType = models.ClientWeb
	case "mobile":
		clientType = models.ClientMobile

	}

	
	// revoke the old token so it can't be used anymore
	if err := RevokeRefreshToken(tokenRecord.TokenText); err != nil {

		log.Printf("Couldn't Revoke the token %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	// Generate new tokens Rotate refresh token (issue a new one)
	tokens, err := issueTokens(c, userID, clientType)

	// store the refresh token 
	newExpireTime := time.Now().Add(24*30*time.Hour)
	_,err = StoreRefreshTokens(userID,tokens.RefreshToken,newExpireTime,clientTypeStr)
	if err != nil{

		log.Printf("couldn't store a new refresh token %v",err)
	}


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
		 c.JSON(http.StatusOK, gin.H{
            "access_token":  tokens.AccessToken,
            "refresh_token": tokens.RefreshToken,
        })
	}

	

}

// Refresh token handler for getting new access tokens
func ExtractRefreshToken(c *gin.Context) (string, error) {

	//Extracting  refresh tokens from mobile app clients
	var refreshToken string
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

func issueTokens(c *gin.Context, userID int, clientType models.ClientType) (*TokenPair, error) {

	accessToken, err := pkg.GenerateAcessToken(userID)
	if err != nil {

		return nil, err
	}
	refreshToken, err, expiresAt := pkg.GenerateRefreshToken(userID)
	if err != nil {

		return nil, err
	}

	_, err = StoreRefreshTokens(userID, refreshToken, expiresAt, string(clientType))

	if err != nil {

		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func StoreRefreshTokens(userID int, refreshToken string, expiresAt time.Time, clientType string) (sql.Result, error) {

	// hashing the token before inserting to a db
	refreshToken = pkg.HashToken(refreshToken)

	query := `INSERT INTO refresh_tokens (user_id,token_text,expires_at,client_type) VALUES($1,$2,$3,$4)`

	result, err := pkg.DB.Exec(query, userID, refreshToken, expiresAt, clientType)
	if err != nil {

		return nil, err
	}

	return result, nil

}

func GetRefreshToken(refreshToken string) (*models.RefreshToken, error) {

	var refreshRecord models.RefreshToken

	// hashing the token because stored tokens are hashed
	hashedrefreshToken := pkg.HashToken(refreshToken)

	query := `SELECT * FROM refresh_tokens where token_text = $1;`

	if err := pkg.DB.QueryRow(query, hashedrefreshToken).Scan(&refreshRecord); err != nil {

		return nil, err
	}

	return &refreshRecord, nil
}

// RevokeRefreshToken revokes the token
func RevokeRefreshToken(refreshToken string) error {

	query := `
	
	UPDATE refresh_tokens SET revoked=TRUE 
	WHERE token_text = $1 revoked=FALSE `

	result, err := pkg.DB.Exec(query, refreshToken)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// detect reuse attempt
	if rowsAffected == 0 {
		// token was already revoked or doesn't exist
		log.Printf("refresh token already revoked or not found")
	}

	return err

}

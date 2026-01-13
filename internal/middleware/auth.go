package middleware

import (
	"log"
	"net/http"
	"strings"
	"github.com/abelmalu/golang-posts/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header in the request header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Stop the request from reaching the handler
			return
		}

		// 2. Check for the "Bearer " prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			log.Printf("Bearer Token not in the request header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// 3. Validate the token
		tokenString := parts[1]
		token, err := pkg.ValidateToken(tokenString) // Your existing function

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}


		// 4. Extract User ID and set it in the Context
		// This allows the CreatePost handler to know WHO is posting
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
            uid, ok := claims["user_id"].(float64)
            if !ok {
                c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
                return
            }
            c.Set("userID", int(uid)) // convert once here
        }

		c.Next() // Token is valid, proceed to the next handler!
	}
}
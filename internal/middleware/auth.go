package middleware

import (
	"log"
	"net/http"
	"strings"
	
	"github.com/abelmalu/golang-posts/pkg"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		// 1. Get the Authorization header
		authHeader := c.GetHeader("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {

			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if tokenStr == "" {

			cookieToken, err := c.Cookie("access_token")
			if err != nil {
				log.Printf("Getting cookie error %v",err)

				
			}
			tokenStr = cookieToken

		}
		if tokenStr == "" {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})

			c.Abort()
			return

		}

		// Validate the token

		tokenClaims, err := pkg.ValidateAccessToken(tokenStr) // Your existing function

		if err != nil {
			log.Printf("invalid token, %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		uid, ok := tokenClaims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
		c.Set("userID", int(uid))
		// convert once here

		c.Next() // Token is valid, proceed to the next handler!
	}
}

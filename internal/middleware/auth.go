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
				log.Printf("Getting cookie error %v", err)

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

		tokenClaims, err := pkg.ValidateAccessToken(tokenStr)

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
		userRole, ok := tokenClaims["userRole"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
		c.Set("userRole", userRole)
		// convert once here

		c.Next() // Token is valid, proceed to the next handler!
	}
}

func AuthorizeRoles(allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		role, ok := c.Get("userRole")
		if !ok {

			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "can't Unauthorized",
			})
			c.Abort()
			return
		}
		userRole := role.(string)

		hasAccess := false
		for _, r := range allowedRoles {
			if r == userRole {
				hasAccess = true
				break

			}

		}
		if !hasAccess {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "You are not authorized for this action",
			})
			c.Abort()

		}
		c.Next()

	}
}

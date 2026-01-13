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
		var tokenStr string 
		// 1. Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		
		if strings.HasPrefix(authHeader,"Bearer "){

			tokenStr = strings.TrimPrefix(authHeader,"Bearer ")
		}

		if tokenStr == ""{

			cookieToken, err := c.Cookie("access_token")
			if err == nil {

				tokenStr = cookieToken
			}

		}
		if tokenStr == ""{

			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"Unauthorized",
			})
		}
		// Validate the token
		
		token, err := pkg.ValidateToken(tokenStr) // Your existing function

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
package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware(logger *platform.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)
				
				errStr :=fmt.Sprintf( "[PANIC] request_id=%s time=%s error=%v\n%s", 
					requestID, 
					time.Now().Format(time.RFC3339), 
					err, 
					debug.Stack())

				logger.Error(errStr)

			
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"type":       "INTERNAL",
					"message":    "An unexpected error occurred",
					"request_id": requestID,
					"timestamp":  time.Now().Unix(),
				})
			}
		}()
		c.Next()
	}
}

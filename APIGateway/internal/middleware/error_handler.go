package middleware

import (
	"log"
	"time"

	"github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only handle if there are errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			requestID := GetRequestID(c)

			var appErr *errors.AppError
			var ok bool

			if appErr, ok = err.(*errors.AppError); !ok {
				// Convert unknown error to internal error
				appErr = errors.NewInternalError("An unexpected error occurred", err)
			}

			// Add RequestID to the error for response
			appErr.RequestID = requestID

			// Log the actual cause if it's an internal error or has a cause
			if appErr.Type == errors.TypeInternal || appErr.Cause != nil {
				log.Printf("[ERROR] request_id=%s type=%s message=%s cause=%v", 
					requestID, appErr.Type, appErr.Message, appErr.Cause)
			}

			// Final API response
			c.AbortWithStatusJSON(appErr.HTTPStatus(), gin.H{
				"type":       appErr.Type,
				"message":    appErr.Message,
				"details":    appErr.Details,
				"request_id": appErr.RequestID,
				"timestamp":  time.Now().Unix(),
			})
		}
	}
}

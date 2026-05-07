package middleware

import (
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

			var appErr *ierrors.AppError
			var ok bool

			if appErr, ok = err.(*ierrors.AppError); !ok {
				// Convert unknown error to internal error
				appErr = ierrors.NewInternalError("An unexpected error occurred", err)
			}

			// Add RequestID to the error for response
			appErr.RequestID = requestID

			

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

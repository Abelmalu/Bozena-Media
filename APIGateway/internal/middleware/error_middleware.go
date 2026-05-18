package middleware

import (
	"time"
	"github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			requestID := GetRequestID(c)

			var appErr *ierrors.AppError
			var ok bool

			if appErr, ok = err.(*ierrors.AppError); !ok {
				appErr = ierrors.NewInternalError("An unexpected error occurred", err)
			}

			appErr.RequestID = requestID

			

			c.AbortWithStatusJSON(appErr.HTTPStatus(), gin.H{
				"type":       appErr.Type,
				"message":    appErr.Message,
				"request_id": appErr.RequestID,
				"timestamp":  time.Now().In(time.FixedZone("EAT",3*60*60)),
			})
		}
	}
}

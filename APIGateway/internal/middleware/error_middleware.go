package middleware

import (
	"github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
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


			
		utils.SendErrorResponse[error](c,appErr,requestID, appErr.HTTPStatus())
		}
	}
}

package middleware

import (
	"errors"
	"time"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware(logger *platform.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		userAgent := c.GetHeader("User-Agent")
		sourceApp := c.GetHeader("source_app")
			requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			c.Abort()
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			c.Abort()
			return

		}
	}
     

		logger.RequestStart(c.Request.Method, c.Request.URL.RequestURI(), userAgent, sourceApp, requestID)

		c.Next()
		statusCode := c.Writer.Status()

		logger.RequestEnd(c.Request.Method, c.Request.URL.RequestURI(), statusCode, time.Since(startTime), requestID)

	}

}

package middleware

import (
	"time"
	"github.com/abelmalu/golang-posts/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)

func RequestLoggerMiddleware(logger *platform.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		userAgent := c.GetHeader("User-Agent")
		sourceApp := c.GetHeader("source_app")
		requestID,err := utils.GetRequestID(c,logger)
		if err != nil{

			c.Abort()
		}
     

		logger.RequestStart(c.Request.Method, c.Request.URL.RequestURI(), userAgent, sourceApp, requestID)

		c.Next()
		statusCode := c.Writer.Status()

		logger.RequestEnd(c.Request.Method, c.Request.URL.RequestURI(), statusCode, time.Since(startTime), requestID)

	}

}

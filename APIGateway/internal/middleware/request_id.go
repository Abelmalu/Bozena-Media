package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"
const RequestIDHeader = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
			requestID := uuid.New().String()
		

		// Set in Gin context for easy access in handlers
		c.Set(string(RequestIDKey), requestID)

		// Set in standard context for propagation to other layers
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// Set in response header
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if val, ok := c.Get(string(RequestIDKey)); ok {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

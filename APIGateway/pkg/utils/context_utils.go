package utils

import (
	"context"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// GetRequestID get the RequestID from the context
func GetRequestID(c *gin.Context) (string, error) {

	requestID, ok := c.Get("request_id")
	if !ok {
		return "", ierrors.ErrRequestIDNotFoundInContext
	}
	requestIDValue, ok := requestID.(string)
	if !ok {

		return "", ierrors.ErrTypeAssertionFailed

	}

	return requestIDValue, nil

}

// GetUserID get the userID from the context
func GetUserID(c *gin.Context) (int, error) {

	var userID interface{}
	var ok bool

	if userID, ok = c.Get("userID"); ok {

		if userIDValue, ok := userID.(int); ok {

			return userIDValue, nil

		} else {

			return 0, ierrors.ErrTypeAssertionFailed

		}

	} else {

		return 0, ierrors.ErrUserIDNotFoundInContext

	}

}

// addToOutgoingContext add data to the outgoing context
func AddToOutgoingContext(c *gin.Context, requestID string) (context.Context, string) {

	clientType := c.GetHeader("X-Client-Type")

	md := metadata.Pairs(
		"x-client-type", clientType,
		"request-id", requestID,
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx, clientType

}

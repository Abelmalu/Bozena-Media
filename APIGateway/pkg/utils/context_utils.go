package utils

import (
	"context"
	"errors"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

//GetRequestID get the RequestID from the context
func GetRequestID(c *gin.Context) (string, error) {

	requestID, ok := c.Get("request_id")
	if !ok {
		return "",ierrors.ErrRequestIDNotFoundInContext
	}
	requestIDValue, ok := requestID.(string)
	if !ok {

	return "",ierrors.ErrTypeAssertionFailed

	}

	return requestIDValue, nil

}

//GetUserID get the userID from the context
func GetUserID(c *gin.Context, logger *platform.Logger) (int, error) {

	var userID interface{}
	var ok bool

	if userID, ok = c.Get("userID"); ok {

		if userIDValue, ok := userID.(int); ok {

			return userIDValue, nil

		} else {

			logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return 0, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil)

		}

	} else {

		logger.Error("couldn't get request ID", zap.Error(errors.New("couldn't find request ID")))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

		return 0, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil)

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


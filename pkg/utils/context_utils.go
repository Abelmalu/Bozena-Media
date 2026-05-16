package utils

import (
	"errors"

	ierrors "github.com/abelmalu/golang-posts/pkg/errors"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetRequestID(c *gin.Context, logger *platform.Logger) (string, error) {

	requestID, ok := c.Get("request_id")
	if !ok {
		logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

		return "", ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil)
	}
	requestIDValue, ok := requestID.(string)
	if !ok {

		logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return "", ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil)

	}

	return requestIDValue, nil

}

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

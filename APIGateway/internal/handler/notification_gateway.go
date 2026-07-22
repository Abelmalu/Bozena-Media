package handler

import (
	"errors"
	"net/http/httputil"
	"net/url"
	"strconv"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NotificationHanlder struct {
	logger *platform.Logger
}

var target, _ = url.Parse("http://localhost:8083")

var proxy = httputil.NewSingleHostReverseProxy(target)

func NewNotificationHandler(logger *platform.Logger) *NotificationHanlder {

	return &NotificationHanlder{

		logger: logger,
	}
}

func (nh *NotificationHanlder) Stream(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			nh.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			nh.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	userID, err := utils.GetUserID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

			nh.logger.Error("couldn't couldn't find userID in the context", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			nh.logger.Error("couldn't assert the user ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}


	userIDStr := strconv.Itoa(userID)

	c.Request.Header.Set("X-Request-ID", requestID)
	c.Request.Header.Set("X-User-ID", userIDStr)

	nh.logger.Info("proxying the request to notification service")
	proxy.ServeHTTP(c.Writer, c.Request)

}

func (nh *NotificationHanlder) GetUserNotifications(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			nh.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			nh.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	userID, err := utils.GetUserID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

			nh.logger.Error("couldn't couldn't find userID in the context", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			nh.logger.Error("couldn't assert the user ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}

	userIDStr := strconv.Itoa(userID)

	c.Request.Header.Set("X-Request-ID", requestID)
	c.Request.Header.Set("X-User-ID", userIDStr)

	nh.logger.Info("proxying the request to notification service")
	proxy.ServeHTTP(c.Writer, c.Request)

}

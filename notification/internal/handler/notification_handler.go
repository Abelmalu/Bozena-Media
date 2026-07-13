package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/abelmalu/golang-posts/notification/internal/broker"
	"github.com/abelmalu/golang-posts/notification/internal/core"
	ierrors "github.com/abelmalu/golang-posts/notification/internal/errors"
	"github.com/abelmalu/golang-posts/notification/pkg"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NotificationHanlder struct {
	logger              *platform.Logger
	notificationService core.NotificationService
	notificationBroker *broker.NotificationBroker
}

func NewNotificationHandler(logger *platform.Logger, service core.NotificationService,notiBroker *broker.NotificationBroker) *NotificationHanlder {

	return &NotificationHanlder{

		logger:              logger,
		notificationService: service,
		notificationBroker: notiBroker,
	}
}

func (notificationHanlder *NotificationHanlder) Stream(c *gin.Context) {

	userID := c.GetHeader("X-User-ID")
	requestID := c.GetHeader("X-Request-ID")
	userIDInt, err := strconv.Atoi(userID)

	if err != nil {
		notificationHanlder.logger.Error("Error", zap.Error(err))

		pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {

		userChan := notificationHanlder.notificationBroker.Register(userIDInt)


		select {
		case <-c.Request.Context().Done():
			return false

		case msg, ok := <-userChan:
			if !ok {
				return false 
			}
			c.SSEvent("notification", msg)
			return true
		}
	})
}

func (notificationHanlder *NotificationHanlder) GetUserNotifications(c *gin.Context) {

	userID := c.GetHeader("X-User-ID")
	requestID := c.GetHeader("X-Request-ID")
	userIDInt, err := strconv.Atoi(userID)

	if err != nil {
		notificationHanlder.logger.Error("Error", zap.Error(err))

		pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)

	if err != nil {

		notificationHanlder.logger.Error("Error", zap.Error(err))

		pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)
		return
	}

	cursor := c.Query("cursor")

	resp, err := notificationHanlder.notificationService.GetUserNotifications(c.Request.Context(), userIDInt, cursor, limit)

	if err != nil {

		notificationHanlder.logger.Error("Error", zap.Error(err))

		var appErr *ierrors.AppError

		if errors.As(err, &appErr) {

			switch appErr.Type {
			case ierrors.TypeDatabase:

				pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)
				return
			default:
				pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)
				return

			}

		}

		pkg.SendErrorResponse[ierrors.AppError](c, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err), requestID, http.StatusInternalServerError)

		return
	}

	pkg.SendSuccessResponse(c, resp, requestID, http.StatusOK)
}

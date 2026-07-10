package handler

import (
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)



type NotificationHanlder struct {
	logger *platform.Logger
}

func NewNotificationHandler(logger *platform.Logger) *NotificationHanlder {

	return &NotificationHanlder{

		logger: logger,
	}
}

func (notificationHanlder *NotificationHanlder) Stream(c *gin.Context) {



}

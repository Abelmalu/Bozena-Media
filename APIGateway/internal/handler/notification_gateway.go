package handler

import (
	"net/http/httputil"
	"net/url"

	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
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

func (notificationHanlder *NotificationHanlder) Stream(c *gin.Context) {

	notificationHanlder.logger.Info("proxying the request to notification service")
	proxy.ServeHTTP(c.Writer, c.Request)

}

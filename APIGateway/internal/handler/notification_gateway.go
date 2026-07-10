package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)

type NotificationHanlder struct {
	logger *platform.Logger
}

var target, _ = url.Parse("http://localhost:8083")

var proxy = httputil.NewSingleHostReverseProxy(target)

func init() {
	// The notification service is mounted at "/" on :8083.
	// Strip the gateway prefix so the upstream sees the path it expects.
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api/notification/stream")
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}

		fmt.Println(req.URL.Path,"the path is")
	}
}

func NewNotificationHandler(logger *platform.Logger) *NotificationHanlder {

	return &NotificationHanlder{

		logger: logger,
	}
}

func (notificationHanlder *NotificationHanlder) Stream(c *gin.Context) {

	notificationHanlder.logger.Info("proxying the request to notification service")
	proxy.ServeHTTP(c.Writer, c.Request)

}

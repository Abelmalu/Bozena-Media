package handler

import (
	"io"
	"time"

	"github.com/abelmalu/golang-posts/notification/internal/core"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)

type NotificationHanlder struct {
	logger *platform.Logger
	notificationService core.NotificationService
}

func NewNotificationHandler(logger *platform.Logger,service core.NotificationService) *NotificationHanlder {

	return &NotificationHanlder{

		logger: logger,
		notificationService: service,
	}
}

func (notificationHanlder *NotificationHanlder) Stream(c *gin.Context) {

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")


	c.Stream(func(w io.Writer) bool {

			time.Sleep(time.Second * 5)
			
			c.SSEvent("this","number is u lucky number")

			// select {
			// // Handle client disconnection
			// case <-c.Request.Context().Done():
			// 	return false

			// // Handle next tick event
			// case t := <-ticker.C:
			// 	c.SSEvent("time-update", map[string]string{
			// 		"time": t.Format(time.RFC3339),
			// 		"user":"hellow",
			// 	})
			// 	return true
			// }

			return true
		})
}

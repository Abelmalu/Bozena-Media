package service

import (
	"context"

	"github.com/abelmalu/golang-posts/notification/internal/core"
)



type NotificationService struct {

	notificationRepo core.NotificationRepository

	
}

func NewNotificationService(notificationRepo core.NotificationRepository) *NotificationService{


	return &NotificationService{}
}

func (notificationService *NotificationService)GetUserNotifications(ctx context.Context,userID int){



}

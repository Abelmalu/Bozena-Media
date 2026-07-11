package core

import "context"



type NotificationService interface {

	GetUserNotifications(ctx context.Context,userID int)



}
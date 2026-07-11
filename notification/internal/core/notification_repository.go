package core

import "context"


type NotificationRepository interface {


	GetUserNotifications(ctx context.Context,userID int)

}
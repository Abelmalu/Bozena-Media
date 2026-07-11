package core

import (
	"context"

	"github.com/abelmalu/golang-posts/notification/internal/dto"
)


type NotificationRepository interface {


	GetUserNotifications(ctx context.Context,userID int , intcursor string,limit int)(*dto.PaginatedResponse,error)
	CreateCacheUser(ctx context.Context, userID int, username, name string) error 
	InsertUserNotification(ctx context.Context,actorID,recipientID int)  error


}
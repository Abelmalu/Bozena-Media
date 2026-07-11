package core

import (
	"context"

	"github.com/abelmalu/golang-posts/notification/internal/dto"
)



type NotificationService interface {

	GetUserNotifications(ctx context.Context, userID int, cursor string, limit int) (*dto.PaginatedResponse, error)
	CreateCacheUser(ctx context.Context, userID int, username, name string) error 
	CreateUserNotification(ctx context.Context,actorID,recipientID int)  error


}
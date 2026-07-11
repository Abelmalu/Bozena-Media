package repository

import (
	"context"
	"database/sql"

)



type NotificationRepository struct {

	DB *sql.DB

	
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository{


	return &NotificationRepository{
		DB:db,
	}
}

func (notificationRepo *NotificationRepository)GetUserNotifications(ctx context.Context,userID int){


	
}

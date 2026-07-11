package service

import (
	"context"

	"github.com/abelmalu/golang-posts/notification/internal/core"
	"github.com/abelmalu/golang-posts/notification/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/notification/internal/errors"
)



type NotificationService struct {

	notificationRepo core.NotificationRepository

	
}

func NewNotificationService(notificationRepo core.NotificationRepository) *NotificationService{


	return &NotificationService{}
}

func (notificationService *NotificationService)GetUserNotifications(ctx context.Context, userID int, cursor string, limit int) (*dto.PaginatedResponse, error){

	resp,err := notificationService.notificationRepo.GetUserNotifications(ctx,userID,cursor,limit)

	if err != nil {

		return nil,err
	}

	return  resp,nil



}


func (notificationService *NotificationService) CreateCacheUser(ctx context.Context,userID int ,username,name string)(error){


	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired,nil,nil)
	}


	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired,nil,nil)
	}
	

  if err := notificationService.notificationRepo.CreateCacheUser(ctx,userID,username,name); err != nil {


	return err


  }

  return nil

}



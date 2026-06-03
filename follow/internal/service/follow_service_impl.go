package service

import (
	"context"

	"github.com/abelmalu/golang-posts/follow/internal/core"
	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/follow/internal/errors"
)



type FollowService struct {

	followRepo core.FollowRepository
}



func NewFollowService(followRepo core.FollowRepository)  *FollowService {


	return &FollowService{
		followRepo: followRepo,
	}
}

func (followService *FollowService) ToggleFollow(ctx context.Context,follow bool,followerID,followingID int)(*dto.FollowResponse,error){

 resp,err := followService.followRepo.ToggleFollow(ctx,follow,followerID,followingID)
 if  err != nil {


		return &dto.FollowResponse{},err
	}

	return &dto.FollowResponse{
		Message: resp,
	},nil



	
}


func (followService *FollowService)	GetUserFollowers(ctx context.Context,followingID,limit int,cursor string)(*dto.PaginatedFollowersResponse,error){

     resp,err := followService.followRepo.GetUserFollowers(ctx,followingID,limit,cursor) 

	 if err != nil {

		return nil,err
	 }

	return resp, nil
}

func (followService *FollowService )	CreateCacheUser(ctx context.Context,userID int ,username,name string)(error){

	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired,nil,nil)
	}


	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired,nil,nil)
	}
	

  if err := followService.followRepo.CreateCacheUser(ctx,userID,username,name); err != nil {


	return err


  }

  return nil



}


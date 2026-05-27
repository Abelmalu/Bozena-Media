package service

import (
	"context"

	"github.com/abelmalu/golang-posts/follow/internal/core"
	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
)



type FollowService struct {

	followRepo core.FollowRepository
}



func NewFollowService(followRepo core.FollowRepository)  *FollowService {


	return &FollowService{
		followRepo: followRepo,
	}
}

func (followService *FollowService) ToggleLike(ctx context.Context,state bool,userID,postID int)(*dto.FollowRequest,error){


	return &dto.FollowRequest{},nil



	
}

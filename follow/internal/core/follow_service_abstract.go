package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
)



type FollowService interface {

		ToggleFollow(ctx context.Context,state bool,followerID,followingID int)(*dto.FollowResponse,error)
     	GetUserFollowers(ctx context.Context,followingID,limit int,cursor string)(*dto.PaginatedFollowersResponse,error)
	    CreateCacheUser(ctx context.Context,userID int ,username,name string)(error)
        GetUserUserFollowings(ctx context.Context, followerId, limit int, cursor string) (*dto.PaginatedFollowingsResponse, error) 
		


}
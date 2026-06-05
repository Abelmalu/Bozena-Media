package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
)


type FollowRepository interface{

	ToggleFollow(ctx context.Context,state bool,followerID,followingID int)(string,error)
	GetUserFollowers(ctx context.Context,followingID,limit int,cursor string)(*dto.PaginatedFollowersResponse,error)
    GetUserUserFollowings(ctx context.Context, followerId, limit int, cursor string) (*dto.PaginatedFollowingsResponse, error) 
	CreateCacheUser(ctx context.Context,userID int ,username,name string)(error)

}
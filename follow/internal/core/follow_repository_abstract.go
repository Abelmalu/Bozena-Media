package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
)


type FollowRepository interface{

	ToggleFollow(ctx context.Context,state bool,followerID,followingID int)(string,error)
	GetUserFollowers(ctx context.Context,followingID int)(*dto.PaginatedResponse,error)

}
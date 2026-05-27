package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
)



type FollowService interface {

		ToggleLike(ctx context.Context,state bool,userID,postID int)(*dto.FollowRequest,error)

}
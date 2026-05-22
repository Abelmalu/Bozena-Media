package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
)


type LikeService interface {

	ToggleLike(ctx context.Context,state bool,userID,postID int)(dto.ToggleLikeResponse,error)


	
}
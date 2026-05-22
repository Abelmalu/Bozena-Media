package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
)


type LikeRepository interface{

	ToggleLike(ctx context.Context,state bool,userID,postID int)(dto.ToggleLikeResponse,error)



}
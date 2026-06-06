package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
)


type LikeRepository interface{

	ToggleLike(ctx context.Context,state bool,userID,postID int)(string,error)
	CreateCacheUser(ctx context.Context,userID int ,username,name string)(error)
	CreateCachePost(ctx context.Context, postID int ,title string)error
	GetPostLikes (ctx context.Context, postID,limit int,cursor string)(*dto.PaginatedPostLikesResponse,error)

}
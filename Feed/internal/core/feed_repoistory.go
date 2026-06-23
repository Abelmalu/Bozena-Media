package core

import (
	"context"

	"github.com/abelmalu/golang-posts/Feed/internal/dto"
)


type FeedRepository interface {

	GetUserFeed(ctx context.Context,cursor string,userID,limit int)(*dto.PaginatedResponse,error)
	CreateCachePost(ctx context.Context, postID int ,title,content string)error
	CreateCacheUser(ctx context.Context,userID int ,username,name string)(error)


}
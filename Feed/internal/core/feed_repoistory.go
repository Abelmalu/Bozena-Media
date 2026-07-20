package core

import (
	"context"

	"github.com/abelmalu/golang-posts/Feed/internal/dto"
)


type FeedRepository interface {

	GetUserFeed(ctx context.Context,cursor string,userID,limit int)(*dto.PaginatedResponse,error)
	CreateCachePost(ctx context.Context, postID int ,title,content,image string)error
	CreateCacheUser(ctx context.Context,userID int ,username,name string)(error)
	InsertFeedEntries(ctx context.Context,followersID []int,postID,ownerID int  ) error 
	IncreaseLikeCount(ctx context.Context, postID int) error 
	DecreaseLikeCount(ctx context.Context, postID int) error 

}
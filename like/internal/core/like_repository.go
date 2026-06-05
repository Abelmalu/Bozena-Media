package core

import (
	"context"
)


type LikeRepository interface{

	ToggleLike(ctx context.Context,state bool,userID,postID int)(string,error)
	CreateCacheUser(ctx context.Context,userID int ,username,name string)(error)
	CreateCachePost(ctx context.Context, postID int ,title string)error

}
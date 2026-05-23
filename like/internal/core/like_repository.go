package core

import (
	"context"
)


type LikeRepository interface{

	ToggleLike(ctx context.Context,state bool,userID,postID int)(string,error)



}
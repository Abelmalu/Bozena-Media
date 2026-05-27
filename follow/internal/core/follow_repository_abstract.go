package core

import (
	"context"

)


type FollowRepository interface{

	ToggleLike(ctx context.Context,state bool,userID,postID int)(string,error)

}
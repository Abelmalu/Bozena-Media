package core

import (
	"context"

)


type FollowRepository interface{

	ToggleFollow(ctx context.Context,state bool,followerID,followingID int)(string,error)

}
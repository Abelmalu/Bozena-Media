package client

import (
	"context"

	"github.com/abelmalu/golang-posts/follow/proto/pb"
	"google.golang.org/grpc"
)


type FollowClient struct {

	client pb.FollowServiceClient
}

func NewFollowClient(conn *grpc.ClientConn) *FollowClient{

	return &FollowClient{
		client:pb.NewFollowServiceClient(conn),
	}
}

func (followClient *FollowClient) ToggleFollow(ctx context.Context, follow bool, followerID,followingID int) (*pb.FollowResponse, error){

	return followClient.client.ToggleFollow(

		ctx,
		&pb.FollowRequest{
			Follow: follow,
			FollowerId: int64(followerID),
			FollowingId: int64(followingID),
		},
	


	)
}


func (pc *FollowClient) GetUserFollowers(ctx context.Context,followingID int,limit int,cursor string)(*pb.GetUserFollowersResponse,error){




	return pc.client.GetUserFollowers(
		ctx,
		&pb.GetUserFollowersRequest{
			FollowingId: int64(followingID),
			Limit: int64(limit),
			Cursor: cursor,
		},
	
	)

	
}


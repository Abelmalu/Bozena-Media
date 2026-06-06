package client

import (
	"context"
	"time"

	"github.com/abelmalu/golang-posts/like/proto/pb"
	"google.golang.org/grpc"
)


type LikeClient struct {

	client pb.LikeServiceClient
}

func NewLikeClient(conn *grpc.ClientConn) *LikeClient {


	return &LikeClient{
		client:pb.NewLikeServiceClient(conn),
	
	}
}


func(likeClient *LikeClient) ToggleLike(ctx context.Context, like bool, opts ...grpc.CallOption) (*pb.LikeResponse, error){
   
	ctx,cancel := context.WithTimeout(ctx,time.Minute * 2)
    defer cancel()

	return  likeClient.client.ToggleLike(
		ctx,
		&pb.LikeRequest{
			State: like,
		},
	
	)
}

func (likeClient *LikeClient) GetPostLikes(ctx context.Context,postID int,limit int,cursor string)(*pb.GetPostLikesResponse,error){


	return likeClient.client.GetPostLikes(
		ctx,
		&pb.GetPostLikesRequest{
			PostId: int64(postID),
			Cursor: cursor,
			Limit: int64(limit),
			
		},
	)
}

 
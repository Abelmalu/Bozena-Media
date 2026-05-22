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
 
package client

import (
	"context"

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


func(likeClient *LikeClient) ToggleLike(ctx context.Context, in *pb.LikeRequest, opts ...grpc.CallOption) (*pb.LikeResponse, error){



	return nil,nil
}
 
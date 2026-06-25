package client

import (
	"context"

	"github.com/abelmalu/golang-posts/Feed/proto/pb"
	"google.golang.org/grpc"
)



type FeedClient struct {

	client pb.FeedServiceClient
}



func NewFeedClient (conn *grpc.ClientConn) *FeedClient{



	return &FeedClient{
		client: pb.NewFeedServiceClient(conn),
	}



}


func (feedClient *FeedClient) GetUserFeed(ctx context.Context,userId,limit int, cursor string)(*pb.GetUserFeedResponse,error) {


	return feedClient.client.GetUserFeed(
		ctx,
		&pb.GetUserFeedRequest{
			UserId: int64(userId),
			Cursor: cursor,
			Limit: int64(limit),

		},

	)
}

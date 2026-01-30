package clients

import (
	"context"
	"time"

	"github.com/abelmalu/golang-posts/post/proto/pb"
	"google.golang.org/grpc"
)



type PostClient struct{

	client pb.PostServiceClient
}

func NewPostClient(conn *grpc.ClientConn) *PostClient{

	return &PostClient{
		client: pb.NewPostServiceClient(conn),
	}
}


func (pc *PostClient) CreatePost(userID int64, title, content string) (*pb.CreatePostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return pc.client.CreatePost(ctx, &pb.CreatePostRequest{
		UserId: userID,
		Title:  title,
		Content: content,
	})
}
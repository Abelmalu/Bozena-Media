package client

import (
	"context"
	"time"

	"github.com/abelmalu/golang-posts/post/proto/pb"
	"google.golang.org/grpc"
)

type PostClient struct {
	client pb.PostServiceClient
}

func NewPostClient(conn *grpc.ClientConn) *PostClient {

	return &PostClient{
		client: pb.NewPostServiceClient(conn),
	}
}

func (pc *PostClient) CreatePost(ctx context.Context, userID int64, title, content string) (*pb.CreatePostResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return pc.client.CreatePost(ctx, &pb.CreatePostRequest{
		UserId:  userID,
		Title:   title,
		Content: content,
	})
}

func (pc *PostClient) ListPosts(ctx context.Context) (*pb.ListPostsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return pc.client.ListPosts(ctx, &pb.ListPostsRequest{})

}

func (pc *PostClient) UpdatePost(ctx context.Context,postID int64, title string,content string) (*pb.UpdatePostResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return pc.client.UpdatePost(ctx, &pb.UpdatePostRequest{

		PostId:  postID,
		Title:   title,
		Content: content,
	})
}

func (pc *PostClient) DeletePost(ctx context.Context,postID int64) (*pb.DeletePostResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return pc.client.DeletePost(
		ctx, &pb.DeletePostRequest{
			PostId: postID},
	)

}

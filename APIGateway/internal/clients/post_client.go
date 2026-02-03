package clients

import (
	"context"
	"fmt"
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
	fmt.Println("the user id %v",userID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return pc.client.CreatePost(ctx, &pb.CreatePostRequest{
		UserId: userID,
		Title:  title,
		Content: content,
	})
}

func (pc *PostClient) ListPosts()(*pb.ListPostsResponse,error){
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()


	return pc.client.ListPosts(ctx,&pb.ListPostsRequest{})


}

func (pc *PostClient) UpdatePost (postID int64, title, content string)(*pb.UpdatePostResponse,error){
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()


	return pc.client.UpdatePost(ctx,&pb.UpdatePostRequest{

		PostId: postID,
		Title: title,
		Content:content,

	})
}

func (pc *PostClient) DeletePost (postID int64)(*pb.DeletePostResponse,error){
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()


	return pc.client.DeletePost(
		ctx,&pb.DeletePostRequest{
			PostId: postID,},
	)

}  
 
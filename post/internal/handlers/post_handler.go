package handlers

import (
	"context"
	"errors"
	"log"

	"github.com/abelmalu/golang-posts/post/internal/core"
	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/abelmalu/golang-posts/post/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// the PostHandler will implement the PostServiceServer
type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service core.PostService
}

func NewPostHandler(service core.PostService) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

func (postHandler *PostHandler) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID: int(req.UserId),
	}
	
	createdPost, err := postHandler.service.CreatePost(ctx, &post)
	if err != nil {

		log.Printf("CreatePost failed: %v", err)

		// Map errors to gRPC status codes
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, "request canceled")
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "timeout")
		}

		if err.Error() == "title required" {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.CreatePostResponse{

		Title:   createdPost.Title,
		Content: createdPost.Content,
	}, nil

}

func (postHandler *PostHandler) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	posts, err := postHandler.service.ListPosts(ctx)
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	pbPosts := make([]*pb.Post, len(posts))

	for i, p := range posts {
		pbPosts[i] = &pb.Post{

			Title:   p.Title,
			Content: p.Content,
			Id:int64(p.ID) ,
		}
	}

	return &pb.ListPostsResponse{
		Posts: pbPosts,
	}, nil
}

func (postHandler *PostHandler) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	postID := int(req.PostId)
	_, err := postHandler.service.UpdatePost(ctx, postID, req.Title, req.Content)

	if err != nil {

		return nil, err
	}

	return &pb.UpdatePostResponse{
		Status:  "Success",
		Message: "Post updated successfully",
	}, nil

}

func (postHandler *PostHandler) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error){


    if err := postHandler.service.DeletePost(ctx,int(req.PostId)); err != nil{

        
        return nil,err
    }
  return &pb.DeletePostResponse{
    Status:"Success",
    Message: "Successfully Deleted a Post",
  },nil



}

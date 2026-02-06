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

func (ph *PostHandler) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
  
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
	}
     createdPost, err := ph.service.CreatePost(ctx,&post)
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

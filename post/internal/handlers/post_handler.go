package handlers

import (
	"context"

	"github.com/abelmalu/golang-posts/post/internal/core"
	"github.com/abelmalu/golang-posts/post/proto/pb"
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
	// extract user_id from context (JWT)
	// insert into DB
	return &pb.CreatePostResponse{
		Status:  "success",
		Message: "Post created",
	}, nil
}

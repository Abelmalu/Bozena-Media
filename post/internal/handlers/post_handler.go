package handlers

import (
	"context"
	"fmt"

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

    fmt.Println(req.Title,req.Content)
    print("nothing here")
    fmt.Println(req.UserId)
    
	


	return &pb.CreatePostResponse{
		Status:  "what the hell no DB",
		Message: "Post created",
	}, nil
}

package main

import (
	"context"
	"net"

	"github.com/abelmalu/golang-posts/post/proto/pb"
	"google.golang.org/grpc"
)
type postServer struct {
  pb.UnimplementedPostServiceServer
}

func (s *postServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
    // extract user_id from context (JWT)
    // insert into DB
    return &pb.CreatePostResponse{
        Status:  "success",
        Message: "Post created",
    }, nil
}

func main() {
   
}



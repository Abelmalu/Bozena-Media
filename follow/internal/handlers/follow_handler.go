package handler

import (
	"context"

	"github.com/abelmalu/golang-posts/follow/internal/core"
	"github.com/abelmalu/golang-posts/follow/proto/pb"
)


type FollowHandler struct {

	followService  core.FollowService

	pb.UnimplementedFollowServiceServer

}


func NewFollowHandler(followService core.FollowService)  *FollowHandler {


	return &FollowHandler{
		followService: followService,
	}
}



func (followHandler *FollowHandler) ToggleFollow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowResponse, error){



	return &pb.FollowResponse{},nil
}
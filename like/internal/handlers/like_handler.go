package handler

import (
	"context"

	"github.com/abelmalu/golang-posts/like/internal/core"
	"github.com/abelmalu/golang-posts/like/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
)



type LikeHandler struct {
	pb.UnimplementedLikeServiceServer


	likeService core.LikeService
	logger *platform.Logger
}


func NewLikeHandler(likeService core.LikeService,logger *platform.Logger) *LikeHandler {


	return &LikeHandler{

		likeService: likeService,
		logger:logger,
	}
}


func (likeHandler *LikeHandler)ToggleLike(ctx context.Context, req *pb.LikeRequest) (*pb.LikeResponse, error){


	return &pb.LikeResponse{},nil

}

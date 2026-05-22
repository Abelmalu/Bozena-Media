package handler

import (
	"context"

	"github.com/abelmalu/golang-posts/like/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)


type LikeService interface {
	ToggleLike(ctx context.Context, in *pb.LikeRequest, opts ...grpc.CallOption) (*pb.LikeResponse, error)


}
type LikeHandler struct {

	logger *platform.Logger
	likeClient LikeService
}


func NewLikeHandler(likeClient LikeService,logger *platform.Logger) *LikeHandler{


	return &LikeHandler{
		likeClient: likeClient,
		logger: logger,
	}
}
func (likeHandler *LikeHandler) ToggleLike(c *gin.Context){



}
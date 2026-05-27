package handler

import (
	"context"
	"net/http"

	"github.com/abelmalu/golang-posts/follow/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)


type FollowService interface {
	ToggleFollow(ctx context.Context, state bool, opts ...grpc.CallOption) (*pb.FollowResponse, error)
}
type FollowHandler struct {
	logger     *platform.Logger
	followClient FollowService
}

func NewFollowHandler(followClient FollowService, logger *platform.Logger) *FollowHandler {

	return &FollowHandler{
		followClient: followClient,
		logger:     logger,
	}
}

func (followHandler *FollowHandler) ToggleFollow(c *gin.Context){


	c.JSON(http.StatusOK,"welcome to the follow page")


}
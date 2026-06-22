package handler

import (
	"context"

	"github.com/abelmalu/golang-posts/Feed/internal/core"
	"github.com/abelmalu/golang-posts/Feed/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
)



type FeedHandler struct {

	feedService core.FeedService
	logger *platform.Logger

}



func NewFeedHandler(service core.FeedService,logger *platform.Logger) *FeedHandler {



	return &FeedHandler{

		feedService: service,
		logger: logger,
	}

}

func (feedHandler *FeedHandler) GetUserFeed(ctx context.Context, req *pb.GetUserFeedRequest) (*pb.GetUserFeedResponse, error) {






return nil,nil

}
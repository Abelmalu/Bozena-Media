package service

import (
	"context"

	"github.com/abelmalu/golang-posts/Feed/internal/core"
	"github.com/abelmalu/golang-posts/Feed/internal/dto"
)



type FeedService struct {


	FeedRepo core.FeedRepository



}

func NewFeedService(feedRepo core.FeedRepository) *FeedService {



	return &FeedService{
		FeedRepo: feedRepo,
	}
}


func (feedService *FeedService)	GetUserFeed(ctx context.Context,userID int)(*dto.PaginatedResponse,error){


	return nil,nil
}

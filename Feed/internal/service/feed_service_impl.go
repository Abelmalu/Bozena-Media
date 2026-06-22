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


func (feedService *FeedService)	GetUserFeed(ctx context.Context,cursor string, userID,limit int)(*dto.PaginatedResponse,error){

		resp,err := feedService.FeedRepo.GetUserFeed(ctx,cursor,userID,limit)

		if err != nil {

			return nil,err
		}

	return resp,nil
}

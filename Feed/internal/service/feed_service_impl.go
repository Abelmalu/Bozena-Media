package service

import "github.com/abelmalu/golang-posts/Feed/internal/core"



type FeedService struct {


	FeedRepo core.FeedRepository



}

func NewFeedService(feedRepo core.FeedRepository) *FeedService {



	return &FeedService{
		FeedRepo: feedRepo,
	}
}



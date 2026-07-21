package service

import (
	"context"

	"github.com/abelmalu/golang-posts/Feed/internal/core"
	"github.com/abelmalu/golang-posts/Feed/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/Feed/internal/errors"
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


func (feedService *FeedService) CreateCachePost(ctx context.Context, postID int, title,content,image string,userID int) error {

	err := feedService.FeedRepo.CreateCachePost(ctx, postID, title,content,image,userID)

	return err

}


func (feedService *FeedService) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if err := feedService.FeedRepo.CreateCacheUser(ctx, userID, username, name); err != nil {

		return err

	}

	return nil

}



func (feedService *FeedService) CreateFeedEntries(ctx context.Context,followersID []int,postID,ownerID int ) error {


	err := feedService.FeedRepo.InsertFeedEntries(ctx,followersID,postID,ownerID)
	

	return err



}





func (feedService *FeedService)IncreaseLikeCount(ctx context.Context,postID int) error {

	if err := feedService.FeedRepo.IncreaseLikeCount(ctx,postID); err != nil {

		return err
	}

	return nil
}
func (feedService *FeedService)	DecreaseLikeCount(ctx context.Context,postID int) error {


	if err := feedService.FeedRepo.DecreaseLikeCount(ctx,postID); err != nil {

		return err
	}

	return nil
}


func (feedService *FeedService) GetCachePosts(ctx context.Context,userID int)(*dto.UserCachePostsResponse,error){

	resp,err := feedService.FeedRepo.GetCachePosts(ctx,userID)

	if err != nil {

		return nil,err
	}

	return resp,nil

}


func (feedService *FeedService) AddFeedEntries(ctx context.Context, feedEntries []*dto.FeedEntry) error{


	if err := feedService.FeedRepo.AddFeedEntries(ctx,feedEntries); err != nil {

		return err
	}

	return nil
}


func (feedService *FeedService)	DeleteFeedEntries(ctx context.Context,userID,ownerID int) error {


	if err := feedService.FeedRepo.DeleteFeedEntries(ctx,userID,ownerID); err != nil {

		return err
	}

	return nil
}

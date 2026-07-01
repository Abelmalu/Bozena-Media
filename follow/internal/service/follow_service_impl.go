package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/follow/internal/core"
	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/follow/internal/errors"
	"github.com/abelmalu/golang-posts/follow/internal/models"
)

type FollowService struct {
	followRepo core.FollowRepository
	kafka      sarama.SyncProducer
}

func NewFollowService(followRepo core.FollowRepository, kafka sarama.SyncProducer) *FollowService {

	return &FollowService{
		followRepo: followRepo,
		kafka:      kafka,
	}
}

func (followService *FollowService) ToggleFollow(ctx context.Context, follow bool, followerID, followingID int) (*dto.FollowResponse, error) {

	resp, err := followService.followRepo.ToggleFollow(ctx, follow, followerID, followingID)
	if err != nil {

		return &dto.FollowResponse{}, err
	}

	var followReturned = models.Follow{

		FollowerID:  followerID,
		FollowingID: followingID,
	}

	var msg *sarama.ProducerMessage

	createdUserByte, err := json.Marshal(followReturned)

	if err != nil {

		return nil, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err)
	}

	switch resp {

	case "followed successfully":
		msg = &sarama.ProducerMessage{
			Topic: "followed",
			Value: sarama.StringEncoder(createdUserByte),
		}
		fmt.Println("followed event sent to kafka")

	case "unfollowed successfully":

		msg = &sarama.ProducerMessage{
			Topic: "unfollowed",
			Value: sarama.StringEncoder(createdUserByte),
		}

	}

	_, _, err = followService.kafka.SendMessage(msg)

	if err != nil {

		return nil, ierrors.NewInternalError(ierrors.ErrorMessage("Kafka Sending Error"), err)
	}


	return &dto.FollowResponse{
		Message: resp,
	}, nil

}

func (followService *FollowService) GetUserFollowers(ctx context.Context, followingID, limit int, cursor string) (*dto.PaginatedFollowersResponse, error) {

	resp, err := followService.followRepo.GetUserFollowers(ctx, followingID, limit, cursor)

	if err != nil {

		return nil, err
	}

	return resp, nil
}

func (followService *FollowService) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if err := followService.followRepo.CreateCacheUser(ctx, userID, username, name); err != nil {

		return err

	}

	return nil

}

func (followService *FollowService) GetUserUserFollowings(ctx context.Context, followerId, limit int, cursor string) (*dto.PaginatedFollowingsResponse, error) {

	resp, err := followService.followRepo.GetUserUserFollowings(ctx, followerId, limit, cursor)

	if err != nil {

		return nil, err
	}

	return resp, nil
}

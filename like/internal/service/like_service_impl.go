package service

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/like/internal/core"
	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/like/internal/errors"
)

type KafkaClient interface {

	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)


}

type LikeService struct {
	likeRepo core.LikeRepository
	kafka    KafkaClient
}

func NewLikeService(likeRepo core.LikeRepository,kafkaCl KafkaClient) *LikeService {

	return &LikeService{
		likeRepo: likeRepo,
		kafka: kafkaCl,
	}
}

func (likeService *LikeService) ToggleLike(ctx context.Context, state bool, userID, postID int) (*dto.ToggleLikeResponse, error) {

	message, err := likeService.likeRepo.ToggleLike(ctx, state, userID, postID)

	if err != nil {

		return nil, err

	}

	var likeCreated struct {
		Id     int `json:"id" db:"id"`
		UserID int `json:"user_id" db:"user_id" validate:"required,gt=0"`
		PostID int `json:"post_id"  db:"post_id" validate:"required,gt=0"`
	}

	likeCreated.UserID = userID
	likeCreated.PostID = postID

	var msg *sarama.ProducerMessage

	createdLikeByte, err := json.Marshal(likeCreated)

	if err != nil {

		return nil, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err)
	}

	switch message {

	case "post liked successfully":

		msg = &sarama.ProducerMessage{
			Topic: "liked",
			Value: sarama.StringEncoder(createdLikeByte),
		}
	default:

		msg = &sarama.ProducerMessage{
			Topic: "unliked",
			Value: sarama.StringEncoder(createdLikeByte),
		}

	}

	_, _, err = likeService.kafka.SendMessage(msg)

	if err != nil {

		return nil, ierrors.NewInternalError(ierrors.ErrorMessage("Kafka Sending Error"), err)
	}


	return &dto.ToggleLikeResponse{
		Message: message,
	}, nil
}

func (likeService *LikeService) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if err := likeService.likeRepo.CreateCacheUser(ctx, userID, username, name); err != nil {

		return err

	}

	return nil

}

func (likeService *LikeService) CreateCachePost(ctx context.Context, postID int, title string) error {

	err := likeService.likeRepo.CreateCachePost(ctx, postID, title)

	return err

}

func (likeService *LikeService) GetPostLikes(ctx context.Context, postID, limit int, cursor string) (*dto.PaginatedPostLikesResponse, error) {

	resp, err := likeService.likeRepo.GetPostLikes(ctx, postID, limit, cursor)

	if err != nil {

		return nil, err
	}

	return resp, nil
}

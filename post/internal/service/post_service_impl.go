package service

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/post/internal/core"
	"github.com/abelmalu/golang-posts/post/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/post/internal/errors"
	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type PostService struct {
	repo        core.PostRepository
	kafka       sarama.SyncProducer
	minioClient *minio.Client
}

func NewPostService(repository core.PostRepository, kafkaClient sarama.SyncProducer, minioClient *minio.Client) *PostService {

	return &PostService{
		repo:        repository,
		kafka:       kafkaClient,
		minioClient: minioClient,
	}
}

func (postService *PostService) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {

	if post.Title == "" {

		return nil, ierrors.NewValidationError(ierrors.MSGTitleIsRrequired, nil, nil)

	}

	createdPost, err := postService.repo.CreatePost(ctx, post)

	if err != nil {

		return nil, err
	}

	createdPostByte, err := json.Marshal(createdPost)

	if err != nil {

		return nil, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err)
	}
	msg := &sarama.ProducerMessage{
		Topic: "postCreated",
		Value: sarama.StringEncoder(createdPostByte),
	}

	// Send the message to Kafka
	_, _, err = postService.kafka.SendMessage(msg)

	if err != nil {
		return nil, ierrors.NewAppError(ierrors.TypeKafka, ierrors.ErrorMessage("Kafka Error production error"), err)
	}

	return createdPost, nil

}
func (postService *PostService) UpdatePost(ctx context.Context, postID int, title, content string) (*models.Post, error) {

	if postID <= 0 {

		return nil, ierrors.NewValidationError(ierrors.MSGPathParamError, nil, nil)
	}

	updatedPost, err := postService.repo.UpdatePost(ctx, postID, title, content)
	if err != nil {

		return nil, err
	}

	return updatedPost, nil

}
func (postService *PostService) DeletePost(ctx context.Context, postID int) error {

	if postID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGTitleIsRrequired, nil, nil)

	}
	if err := postService.repo.DeletePost(ctx, postID); err != nil {

		return err
	}
	return nil
}
func (postService *PostService) ListPosts(ctx context.Context) ([]models.Post, error) {

	posts, err := postService.repo.ListPosts(ctx)

	if err != nil {

		return nil, err
	}

	return posts, nil

}

func (postService *PostService) GetUserPosts(ctx context.Context, UserID, limit int64, cursor string) (*dto.PaginatedResponse, error) {

	resp, err := postService.repo.GetUserPosts(ctx, UserID, limit, cursor)

	if err != nil {

		return nil, err

	}

	return resp, nil
}

func (postService *PostService) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if err := postService.repo.CreateCacheUser(ctx, userID, username, name); err != nil {

		return err

	}

	return nil

}

func (postService *PostService) GenerateUploadURL(ctx context.Context, filename, contentType string, userID int) (string, map[string]string, error) {

	if !dto.AllowedTypes[contentType] {

		return "", nil, ierrors.NewBadRequestError("Invalid file fromat", nil)

	}
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".avif":
	default:
		return "", nil, ierrors.NewBadRequestError("Invalid file fromat", nil)
	}

	objectName := fmt.Sprintf(
		"users/%s%s",
		uuid.New().String(),
		ext,
	)

	policy := minio.NewPostPolicy()

	_ = policy.SetBucket("bozena-media")
	_ = policy.SetKey(objectName)
	_ = policy.SetExpires(time.Now().UTC().Add(time.Minute * 10))
	_ = policy.SetContentType(contentType)
	_ = policy.SetContentLengthRange(1, 5*1024*1024) // 5 megabytes only

	
	url, formData, err := postService.minioClient.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return "", nil, ierrors.NewInternalError("Error generating presigned POST policy", err)
	}

	return url.String(), formData, nil
}

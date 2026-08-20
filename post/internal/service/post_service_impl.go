package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/post/internal/cleanup"
	"github.com/abelmalu/golang-posts/post/internal/core"
	"github.com/abelmalu/golang-posts/post/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/post/internal/errors"
	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type kafkaClient interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)

}

type MinioClient interface {
	PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error)
	PresignedPostPolicy(ctx context.Context, p *minio.PostPolicy) (u *url.URL, formData map[string]string, err error)

}

type PostService struct {
	repo           core.PostRepository
	kafka          kafkaClient
	minioClient    MinioClient
	cleanupService *cleanup.CleanUpService
}

func NewPostService(repository core.PostRepository, kafka kafkaClient, minioClient MinioClient, cleanup *cleanup.CleanUpService) *PostService {

	return &PostService{
		repo:           repository,
		kafka:          kafka,
		minioClient:    minioClient,
		cleanupService: cleanup,
	}
}

func (postService *PostService) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {

	if post.Title == "" {
		postService.cleanupService.DeleteObject("bozena-media", *post.Image)

		return nil, ierrors.NewValidationError(ierrors.MSGTitleIsRrequired, nil, nil)

	}

	createdPost, err := postService.repo.CreatePost(ctx, post)

	if err != nil {

		postService.cleanupService.DeleteObject("bozena-media", *post.Image)

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
func (postService *PostService) UpdatePost(ctx context.Context, postID int, title, content, image string) (*models.Post, error) {

	if postID <= 0 {

		return nil, ierrors.NewValidationError(ierrors.MSGPathParamError, nil, nil)
	}

	updatedPost, err := postService.repo.UpdatePost(ctx, postID, title, content, image)
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
	objectName := ""

	for _, post := range resp.Posts {

		if post.Image != nil {

			objectName = *post.Image

			if objectName == "" {

				continue
			}

			url, err := postService.minioClient.PresignedGetObject(ctx, "bozena-media", objectName, time.Hour, nil)

			if err != nil {

				return nil, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err)
			}

			urlStr := url.String()
			post.Image = &urlStr

		}

	}

	return resp, nil
}

func (postService *PostService) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGUserNotFound, nil, nil)
	}

	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGUsernameIsRequired, nil, nil)
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
		"posts/%s%s",
		uuid.New().String(),
		ext,
	)

	policy := minio.NewPostPolicy()

	_ = policy.SetBucket("bozena-media")
	_ = policy.SetKey(objectName)
	_ = policy.SetExpires(time.Now().UTC().Add(time.Minute * 10))
	_ = policy.SetContentType(contentType)
	_ = policy.SetContentLengthRange(1, 5*1024*1024) // 1-5 megabytes only

	url, formData, err := postService.minioClient.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return "", nil, ierrors.NewInternalError("Error generating presigned POST policy", err)
	}

	return url.String(), formData, nil
}

func (postService *PostService) IncreaseLikeCount(ctx context.Context, postID int) error {

	if err := postService.repo.IncreaseLikeCount(ctx, postID); err != nil {

		return err
	}

	return nil
}
func (postService *PostService) DecreaseLikeCount(ctx context.Context, postID int) error {

	if err := postService.repo.DecreaseLikeCount(ctx, postID); err != nil {

		return err
	}

	return nil
}

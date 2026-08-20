package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/post/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/post/internal/errors"
	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/abelmalu/golang-posts/post/internal/service"
)

type MockPostRepository struct {
	resp *dto.PaginatedResponse
	err  error
}

type MockKafkaClient struct {
	partition int32
	offSet    int32
	err       error
}

func (mk *MockKafkaClient) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {

	return mk.partition, int64(mk.offSet), mk.err
}

func (m *MockPostRepository) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	return nil, m.err
}

func (m *MockPostRepository) UpdatePost(ctx context.Context, ID int, title string, content, image string) (*models.Post, error) {
	return nil, m.err
}

func (m *MockPostRepository) DeletePost(ctx context.Context, postID int) error {
	return m.err
}

func (m *MockPostRepository) ListPosts(ctx context.Context) ([]models.Post, error) {
	return nil, m.err
}

func (m *MockPostRepository) GetUserPosts(ctx context.Context, UserID int64, limit int64, cursor string) (*dto.PaginatedResponse, error) {
	return m.resp, m.err
}

func (m *MockPostRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {
	return m.err
}

func (m *MockPostRepository) IncreaseLikeCount(ctx context.Context, postID int) error {
	return m.err
}

func (m *MockPostRepository) DecreaseLikeCount(ctx context.Context, postID int) error {
	return m.err
}

func TestNotificationService_CreateCacheUser(t *testing.T) {
	tests := []struct {
		name     string
		userID   int
		username string
		userName string
		err      error
		wantErr  bool
		check    func(t *testing.T, err, expected error)
	}{
		{
			name:     "invalid id",
			userID:   0,
			username: "mamo",
			userName: "Mamo",
			wantErr:  true,
			err:      ierrors.NewValidationError(ierrors.MSGUserNotFound, nil, nil),
			check: func(t *testing.T, err, expectedErr error) {

				if err.Error() != expectedErr.Error() {

					t.Fatalf("got different error: %v", err)

				}

			},
		},
		{
			name:   "empty username",
			userID: 1, username: "",
			userName: "Mamo",
			wantErr:  true,
			err:      ierrors.NewValidationError(ierrors.MSGUsernameIsRequired, nil, nil),
			check: func(t *testing.T, err, expectedErr error) {

				if err.Error() != expectedErr.Error() {

					t.Fatalf("got different error: %v", err)

				}

			},
		},
		{
			name:     "repo error",
			userID:   1,
			username: "mamo",
			userName: "Mamo",
			err:      errors.New("db error"),
			wantErr:  true,
			check: func(t *testing.T, err, expectedErr error) {

				if err.Error() != expectedErr.Error() {

					t.Fatalf("got different error: %v", err)

				}

			},
		},
		{
			name:     "success",
			userID:   1,
			username: "alice",
			userName: "Alice",
			wantErr:  false,
			check: func(t *testing.T, err, expected error) {

				if err != nil {

					t.Fatalf("Unexpected error %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockPostRepository{err: tt.err}
			kafka := &MockKafkaClient{

				partition: 1,
				offSet: 2,
				err: nil,
			}
			svc := service.NewPostService(repo,kafka,nil,nil)

			err := svc.CreateCacheUser(context.Background(), tt.userID, tt.username, tt.userName)

			if tt.wantErr && err == nil {
				t.Fatal("expected error got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			tt.check(t, err, tt.err)

		})
	}
}

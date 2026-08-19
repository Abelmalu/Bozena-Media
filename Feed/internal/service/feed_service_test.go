package service_test

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/abelmalu/golang-posts/Feed/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/Feed/internal/errors"
	"github.com/abelmalu/golang-posts/Feed/internal/service"
)

type MockFeedRepository struct {
	resp *dto.PaginatedResponse
	err  error
}

type MockMinioClient struct {
	err error
}

func (m *MockFeedRepository) GetUserFeed(ctx context.Context, cursor string, userID, limit int) (*dto.PaginatedResponse, error) {
	return m.resp, m.err
}

func (m *MockFeedRepository) CreateCachePost(ctx context.Context, postID int, title, content, image string, userID int) error {
	return nil
}

func (m *MockFeedRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {
	return m.err
}

func (m *MockFeedRepository) InsertFeedEntries(ctx context.Context, followersID []int, postID, ownerID int) error {
	return m.err
}

func (m *MockFeedRepository) IncreaseLikeCount(ctx context.Context, postID int) error {
	return nil
}

func (m *MockFeedRepository) DecreaseLikeCount(ctx context.Context, postID int) error {
	return nil
}

func (m *MockFeedRepository) GetCachePosts(ctx context.Context, userID int) (*dto.UserCachePostsResponse, error) {
	return nil, nil
}

func (m *MockFeedRepository) AddFeedEntries(ctx context.Context, feedEntries []*dto.FeedEntry) error {
	return nil
}

func (m *MockFeedRepository) DeleteFeedEntries(ctx context.Context, userID, ownerID int) error {
	return nil
}

func (m *MockMinioClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &url.URL{Scheme: "https", Host: "cdn.example.com", Path: "/" + objectName}, nil
}

func TestFeedService_GetUserFeed(t *testing.T) {
	ctx := context.Background()

	originalImage := "avatar.png"
	nilImage := (*string)(nil)

	tests := []struct {
		name    string
		repo    *MockFeedRepository
		minio   *MockMinioClient
		wantErr bool
		check   func(t *testing.T, resp *dto.PaginatedResponse)
	}{
		{
			name: "success with image rewrite and nil image passthrough",
			repo: &MockFeedRepository{
				resp: &dto.PaginatedResponse{
					UserFeeds: []*dto.UserFeed{
						{ID: 1, PostID: 11, UserName: "alice", Name: "Alice", Image: &originalImage},
						{ID: 2, PostID: 22, UserName: "bob", Name: "Bob", Image: nilImage},
					},
					Cursor:  "next",
					HasNext: true,
				},
			},
			wantErr: false,
			minio:   &MockMinioClient{},
			check: func(t *testing.T, resp *dto.PaginatedResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response")
				}
				if resp.UserFeeds[0].Image == nil {
					t.Fatal("expected first image to be rewritten")
				}
				if got := *resp.UserFeeds[0].Image; got != "https://cdn.example.com/avatar.png" {
					t.Fatalf("unexpected rewritten image: %s", got)
				}
				if resp.UserFeeds[1].Image != nil {
					t.Fatal("expected nil image to stay nil")
				}
				if resp.Cursor != "next" || !resp.HasNext {
					t.Fatal("unexpected pagination fields")
				}
			},
		},
		{
			name: "repository error",
			repo: &MockFeedRepository{
				err: errors.New("db failure"),
			},
			minio:   &MockMinioClient{},
			wantErr: true,
		},
		{
			name: "minio error",
			repo: &MockFeedRepository{
				resp: &dto.PaginatedResponse{
					UserFeeds: []*dto.UserFeed{
						{ID: 1, PostID: 11, UserName: "alice", Name: "Alice", Image: &originalImage},
					},
				},
			},
			minio:   &MockMinioClient{err: errors.New("minio failure")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewFeedService(tt.repo, tt.minio)

			resp, err := svc.GetUserFeed(ctx, "", 42, 10)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			tt.check(t, resp)
		})
	}
}


func TestFeedService_CreateCacheUser(t *testing.T) {
	tests := []struct {
		name     string
		userID   int
		username string
		userName string
		err     error
		wantErr bool
		check   func(t *testing.T, err, expected error)
	}{
		{
			name:     "invalid id",
			userID:   0,
			username: "alice",
			userName: "Alice",
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
			userName: "Alice",
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
			username: "alice",
			userName: "Alice",
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
			repo := &MockFeedRepository{err: tt.err}
			svc := service.NewFeedService(repo, &MockMinioClient{})

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

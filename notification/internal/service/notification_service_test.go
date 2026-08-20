package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/abelmalu/golang-posts/notification/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/notification/internal/errors"
	"github.com/abelmalu/golang-posts/notification/internal/service"
)

type MockNotificationRepository struct {
	resp *dto.PaginatedResponse
	err  error
}

func (m *MockNotificationRepository) GetUserNotifications(ctx context.Context, userID int, intcursor string, limit int) (*dto.PaginatedResponse, error) {
	return m.resp, m.err
}

func (m *MockNotificationRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {
	return m.err
}

func (m *MockNotificationRepository) InsertUserNotification(ctx context.Context, actorID, recipientID int) error {
	return m.err
}

func (m *MockNotificationRepository) GetUser(ctx context.Context, userID int) (*dto.User, error) {
	return nil, m.err
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
			repo := &MockNotificationRepository{err: tt.err}
			svc := service.NewNotificationService(repo)

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

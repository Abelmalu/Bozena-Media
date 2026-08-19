package service_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/Chat/internal/errors"
	"github.com/abelmalu/golang-posts/Chat/internal/models"
	"github.com/abelmalu/golang-posts/Chat/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc/metadata"
)

type MockChatRepository struct {
	chats []*models.Conversation
}

var conversations = []*models.Conversation{
	{
		ID: bson.NewObjectID(),
		Participants: []*models.Participant{
			{UserID: 101, Username: "alice_w"},
			{UserID: 102, Username: "bob_k"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID: bson.NewObjectID(),
		Participants: []*models.Participant{
			{UserID: 103, Username: "charlie_m"},
			{UserID: 104, Username: "dana_r"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID: bson.NewObjectID(),
		Participants: []*models.Participant{
			{UserID: 105, Username: "evan_b"},
			{UserID: 106, Username: "fiona_g"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID: bson.NewObjectID(),
		Participants: []*models.Participant{
			{UserID: 107, Username: "george_t"},
			{UserID: 108, Username: "helen_p"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID: bson.NewObjectID(),
		Participants: []*models.Participant{
			{UserID: 109, Username: "ian_f"},
			{UserID: 110, Username: "julia_v"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
var ctx = metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-client-type", "web"))

type MockMinioClient struct{}

func (m *MockChatRepository) InserMessages(ctx context.Context, senderID, receiverID int, chatID bson.ObjectID, message string) error {

	return nil
}

func (m *MockChatRepository) GetChatBetweenUsers(ctx context.Context, senderID, receiverID int) (*models.Conversation, error) {

	return nil, nil
}

func (m *MockChatRepository) GetUserChats(ctx context.Context, userID, limit int, lastSeenID bson.ObjectID) (*dto.UserChatsResponse, error) {

	for _, chat := range m.chats {

		for _, participant := range chat.Participants {

			if participant.UserID == userID {

				return &dto.UserChatsResponse{

					Chats: []*models.Conversation{chat},
				}, nil

			}
		}
	}

	return nil, ierrors.NewNotFoundError(ierrors.ErrorMessage("User doesn't have any chats"),nil)

}

func (m *MockChatRepository) InsertCacheUser(ctx context.Context, ID int, Username, Name, Avatar string) error {

	return nil
}

func (m *MockChatRepository) UpdateUserAvatar(ctx context.Context, userID int, Avatar string) error {

	return nil
}

func (m *MockChatRepository) GetChatMessages(ctx context.Context, chatID, cursor bson.ObjectID, limit int) (*dto.ChatMessagesResponse, error) {

	return nil, nil
}

func (m *MockMinioClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {

	return nil, nil
}

func TestChatService_GetUsersChat_Success(t *testing.T) {

	mr := &MockChatRepository{

		chats: conversations,
	}

	sv := service.NewChatService(mr, &MockMinioClient{})

	_, err := sv.GetUserChats(ctx, mr.chats[0].Participants[0].UserID, 1, mr.chats[0].ID)

	if err != nil {

		t.Fatalf("unexpected error: %v", err)
	}

}

func TestChatService_GetUsersChat_Error(t *testing.T) {

	mr := &MockChatRepository{

		chats: conversations,
	}

	sv := service.NewChatService(mr, &MockMinioClient{})

	_, err := sv.GetUserChats(ctx, 2, 1, bson.NewObjectID())

	if err == nil {

		t.Fatalf("expecting error got nil")
	}

}

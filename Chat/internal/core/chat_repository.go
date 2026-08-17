package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
	"github.com/abelmalu/golang-posts/Chat/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatRespository interface {
	InserMessages(ctx context.Context, senderID, receiverID int, chatID bson.ObjectID, message string) error
	GetChatBetweenUsers(ctx context.Context, senderID, receiverID int) (*models.Conversation, error)
	GetUserChats(ctx context.Context, userID, limit int, lastSeenID bson.ObjectID) (*dto.UserChatsResponse, error)
	InsertCacheUser(ctx context.Context, ID int, Username, Name, Avatar string) error
	UpdateUserAvatar(ctx context.Context, userID int, Avatar string) error
	GetChatMessages(ctx context.Context, chatID, cursor bson.ObjectID, limit int) (*dto.ChatMessagesResponse, error)
}

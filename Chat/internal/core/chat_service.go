package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatService interface {
	SendMessage(ctx context.Context, senderID, receiverID int, message string) error
	GetUserChats(ctx context.Context, userID, limit int, lastSeenID bson.ObjectID) (*dto.UserChatsResponse, error)
	CreateCacheUser(ctx context.Context, ID int, Username, Name, Avatar string) error
	UpdateUserAvatar(ctx context.Context, userID int, Avatar string) error
	GetChatMessages(ctx context.Context,chatID bson.ObjectID) (*dto.ChatMessagesResponse,error)
}

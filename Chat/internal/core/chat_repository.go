package core

import (
	"context"

	"github.com/abelmalu/golang-posts/Chat/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatRespository interface {
	InserMessages(ctx context.Context, senderID, receiverID int, chatID bson.ObjectID, message string) error
	GetChatBetweenUsers(ctx context.Context, senderID, receiverID int) (*models.Conversation, error)
	GetUserChats(ctx context.Context,userID int)([] *models.Conversation,error)
}

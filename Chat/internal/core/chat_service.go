package core

import (
	"context"

	"github.com/abelmalu/golang-posts/Chat/internal/models"
)

type ChatService interface {
	SendMessage(ctx context.Context, senderID, receiverID int, message string) error
	GetUserChats(ctx context.Context, userID int) ([]*models.Conversation, error)
}

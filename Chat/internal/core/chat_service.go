package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
)

type ChatService interface {
	CreateMessages(ctx context.Context, senderID, receiverID int) *dto.MessageResponse
}

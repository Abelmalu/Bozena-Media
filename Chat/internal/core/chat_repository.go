package core

import (
	"context"

	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
)

type ChatRespository interface {
	InserMessages(ctx context.Context, senderID, receiverID int) *dto.MessageResponse
}

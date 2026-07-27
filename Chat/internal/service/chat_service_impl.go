package service

import (
	"context"

	"github.com/abelmalu/golang-posts/Chat/internal/core"
	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
)

type ChatService struct {
	repo core.ChatRespository
}

func NewChatService(r core.ChatRespository) *ChatService {

	return &ChatService{
		repo: r,
	}
}

func (cs *ChatService) CreateMessages(ctx context.Context, senderID, receiverID int) *dto.MessageResponse {

	return nil
}

package service

import (
	"context"

	"github.com/abelmalu/golang-posts/Chat/internal/core"
)

type ChatService struct {
	repo core.ChatRespository
}

func NewChatService(r core.ChatRespository) *ChatService {

	return &ChatService{
		repo: r,
	}
}

func (cs *ChatService) SendMessage(ctx context.Context, senderID, receiverID int,message string) error{


	return nil 
	 
}

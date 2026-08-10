package service

import (
	"context"

	"github.com/abelmalu/golang-posts/Chat/internal/core"
	"github.com/abelmalu/golang-posts/Chat/internal/models"
)

type ChatService struct {
	repo core.ChatRespository
}

func NewChatService(r core.ChatRespository) *ChatService {

	return &ChatService{
		repo: r,
	}
}

func (cs *ChatService) SendMessage(ctx context.Context, senderID, receiverID int, message string) error {

	chat, err := cs.repo.GetChatBetweenUsers(ctx, senderID, receiverID)
	if err != nil {

		return err
	}

	if err := cs.repo.InserMessages(ctx, senderID, receiverID, chat.ID, message); err != nil {

		return err
	}

	return nil

}

func (cs *ChatService) GetUserChats(ctx context.Context, userID int) ([]*models.Conversation, error) {
	
	resp,err := cs.repo.GetUserChats(ctx,userID)

	if err != nil {

		return nil,err
	}

	return resp, nil
}

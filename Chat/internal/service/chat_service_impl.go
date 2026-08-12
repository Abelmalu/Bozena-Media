package service

import (
	"context"

	"github.com/abelmalu/golang-posts/Chat/internal/core"
	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func (cs *ChatService) GetUserChats(ctx context.Context, userID, limit int, lastSeenID bson.ObjectID) (*dto.UserChatsResponse, error) {

	resp, err := cs.repo.GetUserChats(ctx, userID,limit,lastSeenID)

	if err != nil {

		return nil, err
	}

	return resp, nil
}


func (cs *ChatService) 	CreateCacheUser(ctx context.Context, userID int, Username, Name,Avatar string) error {

	if err := cs.repo.InsertCacheUser(ctx,userID,Username,Name,Avatar); err != nil {

		return err



	}

	return nil

}


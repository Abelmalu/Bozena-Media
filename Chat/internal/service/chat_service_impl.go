package service

import (
	"context"
	"time"

	"github.com/abelmalu/golang-posts/Chat/internal/core"
	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/Chat/internal/errors"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatService struct {
	repo  core.ChatRespository
	minio *minio.Client
}

func NewChatService(r core.ChatRespository, minioClient *minio.Client) *ChatService {

	return &ChatService{
		repo:  r,
		minio: minioClient,
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

	resp, err := cs.repo.GetUserChats(ctx, userID, limit, lastSeenID)

	if err != nil {

		return nil, err
	}

	for _, chat := range resp.Chats {

		for _, user := range chat.Participants {

			objectName := user.Avatar

			if objectName == "" {

				continue
			}
			url, err := cs.minio.PresignedGetObject(ctx, "bozena-media", objectName, time.Hour, nil)

			if err != nil {

				return nil, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err)
			}

			urlStr := url.String()
			user.Avatar = urlStr

		}
	}

	return resp, nil
}

func (cs *ChatService) CreateCacheUser(ctx context.Context, userID int, Username, Name, Avatar string) error {

	if err := cs.repo.InsertCacheUser(ctx, userID, Username, Name, Avatar); err != nil {

		return err

	}

	return nil

}

func (cs *ChatService) UpdateUserAvatar(ctx context.Context, userID int, Avatar string) error {

	err := cs.repo.UpdateUserAvatar(ctx, userID, Avatar)

	return err
}

func (cs *ChatService) GetChatMessages(ctx context.Context, chatID, cursor bson.ObjectID, limit int) (*dto.ChatMessagesResponse, error) {

	resp, err := cs.repo.GetChatMessages(ctx, chatID,cursor,limit)

	if err != nil {

		return nil, err
	}

	return resp, err

}

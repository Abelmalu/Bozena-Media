package repository

import (
	"context"
	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatRepository struct {
	DB *mongo.Database
}

func NewChatRepository(DB *mongo.Database) *ChatRepository {

	return &ChatRepository{
		DB: DB,
	}
}

func (cr *ChatRepository) InserMessages(ctx context.Context, senderID, receiverID int) *dto.MessageResponse {


	return &dto.MessageResponse{


	}

}

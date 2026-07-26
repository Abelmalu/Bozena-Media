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

func (cs *ChatRepository) CreateMessages(ctx context.Context, senderID, receiverID int) *dto.MessageResponse {

}

package repository

import (
	"context"

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

func (cr *ChatRepository) InserMessages(ctx context.Context, senderID, receiverID int) {

}

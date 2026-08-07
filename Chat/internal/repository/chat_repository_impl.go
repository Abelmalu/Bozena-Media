package repository

import (
	"context"
	"time"

	"github.com/abelmalu/golang-posts/Chat/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func (cr *ChatRepository) InserMessages(ctx context.Context, senderID, receiverID int, chatID bson.ObjectID, message string) error {

	chatMessage := models.Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  message,
	}
	_, err := cr.DB.Collection("Messages").InsertOne(ctx, chatMessage)

	if err != nil {
		return err
	}
	return nil
}

func (cr *ChatRepository) FindChat(ctx context.Context, senderID, receiverID int) (*models.Conversation, error) {

	filter := bson.M{
		"participants": bson.M{
			"$all": []bson.M{
				{
					"$elemMatch": bson.M{
						"userID": senderID,
					},
				},
				{
					"$elemMatch": bson.M{
						"userID": receiverID,
					},
				},
			},
		},
	}

	var result models.Conversation
	err := cr.DB.Collection("Chats").FindOne(ctx, filter).Decode(&result)

	if err != nil {

		if err == mongo.ErrNoDocuments {

			now := time.Now()

			conversation := models.Conversation{
				ID: bson.NewObjectID(),

				Participants: []models.Participant{
					{
						UserID: senderID,
					},
					{

						UserID: receiverID,
					},
				},

				CreatedAt: now,
				UpdatedAt: now,
			}

			_, err = cr.DB.Collection("Chats").InsertOne(ctx, conversation)

			if err != nil {
				return nil, err
			}
			return &conversation, nil

		}
	}

	return &result, nil
}

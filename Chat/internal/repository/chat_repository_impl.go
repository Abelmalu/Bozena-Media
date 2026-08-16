package repository

import (
	"context"
	"log"
	"time"

	dto "github.com/abelmalu/golang-posts/Chat/internal/dtos"
	"github.com/abelmalu/golang-posts/Chat/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	if err := cr.UpdateLastChatMessage(ctx, senderID, chatID, message); err != nil {

		return err

	}

	return nil
}

func (cr *ChatRepository) GetChatBetweenUsers(ctx context.Context, senderID, receiverID int) (*models.Conversation, error) {

	filter := bson.M{
		"participants": bson.M{
			"$all": []bson.M{
				{
					"$elemMatch": bson.M{
						"userId": senderID,
					},
				},
				{
					"$elemMatch": bson.M{
						"userId": receiverID,
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

				Participants: []*models.Participant{
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

func (cr *ChatRepository) GetUserChats(ctx context.Context, userID, limit int, lastSeenID bson.ObjectID) (*dto.UserChatsResponse, error) {

	var hasNext bool
	var after bson.ObjectID
	var chats []*models.Conversation
	filter := bson.M{
		"participants.userId": userID,
	}
	if lastSeenID != bson.NilObjectID {
		filter["_id"] = bson.M{"$gt": lastSeenID}
	}

	findOptions := options.Find().SetLimit(int64(limit + 1))

	cursor, err := cr.DB.Collection("Chats").Find(ctx, filter, findOptions)
	if err != nil {

		return nil, err
	}

	for cursor.Next(ctx) {

		var chat models.Conversation

		if err := cursor.Decode(&chat); err != nil {

			return nil, err
		}

		if len(chats) == limit {

			hasNext = true
			after = chat.ID

			break
		}

		chats = append(chats, &chat)

	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &dto.UserChatsResponse{

		Chats:   chats,
		Cursor:  after,
		HasNext: hasNext,
	}, nil
}

func (cr *ChatRepository) UpdateLastChatMessage(ctx context.Context, senderID int, chatID bson.ObjectID, message string) error {

	update := bson.M{
		"$set": bson.M{
			"lastMessage": struct {
				SenderID int    `bson:"senderId"`
				Text     string `bson:"text"`
			}{
				SenderID: senderID,
				Text:     message,
			},
		},
	}

	result, err := cr.DB.Collection("Chats").UpdateByID(ctx, chatID, update)

	if err != nil {

		log.Println("the error", err)

		return err
	}

	if result.ModifiedCount == 0 {

		return nil
	}

	return err

}

func (cr *ChatRepository) InsertCacheUser(ctx context.Context, userID int, Username, Name, Avatar string) error {

	user := models.User{

		UserID:   userID,
		Name:     Name,
		Username: Username,
		Avatar:   Avatar,
	}

	_, err := cr.DB.Collection("users_cache").InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (cr *ChatRepository) UpdateCacheUser(ctx context.Context, userID int, Avatar string) error{


	filter := bson.M{
		"userId":userID,
	}

	update := bson.M{

		"$set":bson.M{
			"avatar":Avatar,

		},
	}

	_, err := cr.DB.Collection("users_cache").UpdateOne(ctx,filter,update)

	if err != nil {
		return err
	}

	return nil

}

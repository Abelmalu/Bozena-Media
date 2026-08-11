package dto

import (
	"github.com/abelmalu/golang-posts/Chat/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MessageRequest struct {
	Message    string `json:"message"`
	ReceiverID int    `json:"receiver_id"`
}

type UserChatsResponse struct {
	Chats   []*models.Conversation `json:"chats"`
	HasNext bool                   `json:"has_next"`
	Cursor  bson.ObjectID                 `json:"cursor"`
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Conversation struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Participants []Participant `bson:"participants" json:"participants"`
	LastMessage  LastMessage   `bson:"lastMessage,omitempty" json:"lastMessage"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Participant struct {
	UserID   int    `bson:"userId" json:"userId"`
	Username string `bson:"username" json:"username"`
	Avatar   string `bson:"avatar" json:"avatar"`
}

type LastMessage struct {
	Text      string    `bson:"text" json:"text"`
	SenderID  int       `bson:"senderId" json:"senderId"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}




type Message struct {
	ID             bson.ObjectID `json:"id" bson:"_id"`
	ChatID bson.ObjectID `json:"chatID" bson:"chatID"`
	SenderID       int                `json:"senderId" bson:"senderId"`
	Content        string             `json:"content" bson:"content"`
	Status         string             `json:"status" bson:"status"`
	CreatedAt      string             `json:"createdAt" bson:"createdAt"`
}


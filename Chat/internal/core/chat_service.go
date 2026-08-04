package core

import (
	"context"
)

type ChatService interface {
	SendMessage(ctx context.Context, senderID, receiverID int,message string)
}

package core

import (
	"context"
)

type ChatRespository interface {
	InserMessages(ctx context.Context, senderID, receiverID int)
}

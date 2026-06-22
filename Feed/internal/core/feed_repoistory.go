package core

import (
	"context"

	"github.com/abelmalu/golang-posts/Feed/internal/dto"
)


type FeedRepository interface {

	GetUserFeed(ctx context.Context,userID int)(*dto.PaginatedResponse,error)


}
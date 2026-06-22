package core

import (
	"context"

	"github.com/abelmalu/golang-posts/Feed/internal/dto"
)




type FeedService interface {


GetUserFeed(ctx context.Context,cursor string,userID,limit int)(*dto.PaginatedResponse,error)


}
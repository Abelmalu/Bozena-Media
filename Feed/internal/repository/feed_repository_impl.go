package repository

import (
	"context"
	"database/sql"

	"github.com/abelmalu/golang-posts/Feed/internal/dto"
)



type FeedRepository struct {

	DB *sql.DB
}





func NewFeedRepository(db *sql.DB) *FeedRepository {



	return &FeedRepository{
		DB: db,
	}
}


func (feedRepo *FeedRepository)	GetUserFeed(ctx context.Context,userID int)(*dto.PaginatedResponse,error){


	return nil,nil
}

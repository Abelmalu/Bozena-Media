package repository

import (
	"context"
	"database/sql"
)


type FollowRepository struct {

	DB *sql.DB
}


func NewFollowRepository(DB *sql.DB)  *FollowRepository {


	return &FollowRepository{
		DB: DB,
	}
}

func (followRepository *FollowRepository) ToggleLike(ctx context.Context,state bool,userID,postID int)(string,error){


	return "",nil
}

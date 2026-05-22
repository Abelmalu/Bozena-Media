package repository

import (
	"context"
	"database/sql"

	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
)



type LikeRespository struct{

	DB *sql.DB
}

func NewLikeRepository(db *sql.DB)(*LikeRespository){

	return &LikeRespository{
		DB: db,
	}
}

func (likeRespository *LikeRespository) ToggleLike(ctx context.Context,state bool,userID,postID int)(dto.ToggleLikeResponse,error){
	

return dto.ToggleLikeResponse{},nil
}
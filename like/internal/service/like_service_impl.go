package service

import (
	"context"

	"github.com/abelmalu/golang-posts/like/internal/core"
	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
)

type LikeService struct {

	likeRepo core.LikeRepository
}

func NewLikeRepository(likeRepo core.LikeRepository ) *LikeService{


	return &LikeService{
		likeRepo: likeRepo,
	}
}


func (likeService *LikeService)	ToggleLike(ctx context.Context,state bool,userID,postID int)(dto.ToggleLikeResponse,error){


	return dto.ToggleLikeResponse{},nil
}

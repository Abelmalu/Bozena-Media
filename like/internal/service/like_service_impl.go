package service

import (
	"context"

	"github.com/abelmalu/golang-posts/like/internal/core"
	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/like/internal/errors"
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
	
	message,err := likeService.likeRepo.ToggleLike(ctx,state,userID,postID)

	 if err != nil {

		return  dto.ToggleLikeResponse{},err
	 }

	return dto.ToggleLikeResponse{
		Message: message,
	},nil
}



func (likeService *LikeService )	CreateCacheUser(ctx context.Context,userID int ,username,name string)(error){

	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired,nil,nil)
	}


	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired,nil,nil)
	}
	

  if err := likeService.likeRepo.CreateCacheUser(ctx,userID,username,name); err != nil {


	return err


  }

  return nil



}



	func (likeService *LikeService) CreateCachePost(ctx context.Context, postID int,title string)error{

			err := likeService.likeRepo.CreateCachePost(ctx , postID ,title)

			return err


	}

package service

import (
	"context"

	"github.com/abelmalu/golang-posts/post/internal/core"
	"github.com/abelmalu/golang-posts/post/internal/dto"
	"github.com/abelmalu/golang-posts/post/internal/errors"
	"github.com/abelmalu/golang-posts/post/internal/models"
)

type PostService struct {
	repo core.PostRepository
}

func NewPostService(repository core.PostRepository) *PostService {

	return &PostService{
		repo: repository,
	}
}

func (postService *PostService) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {

	if post.Title == ""{

		return nil,ierrors.NewValidationError(ierrors.MSGTitleIsRrequired,nil,nil)


	}

	createdPost, err := postService.repo.CreatePost(ctx, post)

	if err != nil {

		return nil, err
	}

	return createdPost, nil

}
func (postService *PostService) UpdatePost(ctx context.Context, postID int, title, content string) (*models.Post, error) {

	if postID <= 0 {

		return nil, ierrors.NewValidationError(ierrors.MSGPathParamError, nil, nil)
	}

	updatedPost, err := postService.repo.UpdatePost(ctx, postID, title, content)
	if err != nil {

		return nil, err
	}

	return updatedPost, nil

}
func (postService *PostService) DeletePost(ctx context.Context, postID int) error {

	if postID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGTitleIsRrequired, nil, nil)

	}
	if err := postService.repo.DeletePost(ctx, postID); err != nil {

		return err
	}
	return nil
}
func (postService *PostService) ListPosts(ctx context.Context) ([]models.Post, error) {

	posts, err := postService.repo.ListPosts(ctx)

	if err != nil {

		return nil, err
	}

	return posts, nil

}


func (postService *PostService) GetUserPosts(ctx context.Context,UserID,limit int64,cursor string)(*dto.PaginatedResponse, error){

resp,err := postService.repo.GetUserPosts(ctx,UserID,limit,cursor)

if err != nil {

	return nil,err

}

	return resp,nil
}




func (postService *PostService) CreateCacheUser(ctx context.Context,userID int ,username,name string)(error){


	if userID <= 0 {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired,nil,nil)
	}


	if username == "" {

		return ierrors.NewValidationError(ierrors.MSGNameIsRequired,nil,nil)
	}
	

  if err := postService.repo.CreateCacheUser(ctx,userID,username,name); err != nil {


	return err


  }

  return nil

}

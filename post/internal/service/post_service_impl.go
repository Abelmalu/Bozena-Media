package service

import (
	"context"

	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/abelmalu/golang-posts/post/internal/core"
)

type PostService struct {
	repo core.PostRepository
}


func NewPostService(repository core.PostRepository) *PostService{


	return &PostService{
		repo:repository,
	}
}


func (ps *PostService) CreatePost(post *models.Post)(*models.Post,error){

	panic("")
}
func (ps *PostService) UpdatePost(postID string)(*models.Post,error){

	panic("")
}
func (ps *PostService) DeletePost(postID string)(error){

	panic("")
}
func (ps *PostService) ListPosts(ctx context.Context)([]models.Post,error){

	panic("")
}

	
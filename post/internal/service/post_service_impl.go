package service

import (
	"context"
	"log"

	"github.com/abelmalu/golang-posts/post/internal/core"
	"github.com/abelmalu/golang-posts/post/internal/models"
)

type PostService struct {
	repo core.PostRepository
}


func NewPostService(repository core.PostRepository) *PostService{


	return &PostService{
		repo:repository,
	}
}


func (postService *PostService) CreatePost(ctx context.Context,post *models.Post)(*models.Post,error){

	createdPost,err := postService.repo.CreatePost(ctx,post)

	if err != nil{

		return nil,err
	}

	return createdPost,nil

	
}
func (ps *PostService) UpdatePost(postID string)(*models.Post,error){

	panic("")
}
func (ps *PostService) DeletePost(postID string)(error){

	panic("")
}
func (postService *PostService) ListPosts(ctx context.Context)([]models.Post,error){
   
	posts,err := postService.repo.ListPosts(ctx)

	if err != nil{

		log.Printf("the error is %v",err)
		return nil,err
	}

	return posts,nil
	
}

	
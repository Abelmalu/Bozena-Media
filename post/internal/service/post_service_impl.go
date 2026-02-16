package service

import (
	"context"
	"errors"
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
func (postService *PostService) UpdatePost(ctx context.Context,postID int,title,content string) (*models.Post, error){

	if postID <= 0{

		return nil, errors.New("post not found")
	}

	updatedPost,err := postService.repo.UpdatePost(ctx,postID,title,content)
	if err != nil{

		return nil,err
	}

	return updatedPost,nil

}
func (postService *PostService) DeletePost(ctx context.Context,postID int) error{

	
	if postID <= 0{

		return  errors.New("post not found")
	}
	if err := postService.repo.DeletePost(ctx,postID); err != nil{

		return err
	}
	return nil
}
func (postService *PostService) ListPosts(ctx context.Context)([]models.Post,error){
   
	posts,err := postService.repo.ListPosts(ctx)

	log.Printf("posts inside of service")
	log.Print(posts[1].Content)
	

	if err != nil{

		log.Printf("the error is %v",err)
		return nil,err
	}

	return posts,nil
	
}

	
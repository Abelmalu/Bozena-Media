package core

import (
	"context"
	"github.com/abelmalu/golang-posts/post/internal/models"

)


type PostService interface {
    //CreatePost creates a post 
	CreatePost(ctx context.Context,post *models.Post)(*models.Post,error)
	UpdatePost(ctx context.Context,ID int  )(error)
	DeletePost(ctx context.Context, ID int )(error)
	ListPosts(ctx context.Context)([]models.Post,error)



}
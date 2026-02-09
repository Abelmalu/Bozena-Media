package core

import (
	"context"

	"github.com/abelmalu/golang-posts/post/internal/models"
)

type PostService interface {
	CreatePost(ctx context.Context,post *models.Post) (*models.Post, error)
	UpdatePost(ctx context.Context,postID int,title,content string) (*models.Post, error)
	DeletePost(ctx context.Context,postID int) error
	ListPosts(ctx context.Context) ([]models.Post, error)
	
}

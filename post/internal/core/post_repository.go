package core

import (
	"context"

	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/abelmalu/golang-posts/post/internal/repository"
)

type PostRepository interface {

	CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	UpdatePost(ctx context.Context, ID int, title string,content string)  (*models.Post, error)
    DeletePost(ctx context.Context,postID int)(error)
	ListPosts(ctx context.Context) ([]models.Post, error)
	GetUserPosts(ctx context.Context,UserID int64,limit int64)(*repository.PaginatedResponse, error)	
}
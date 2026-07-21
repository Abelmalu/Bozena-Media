package core

import (
	"context"

	"github.com/abelmalu/golang-posts/post/internal/dto"
	"github.com/abelmalu/golang-posts/post/internal/models"
)

type PostService interface {
	CreatePost(ctx context.Context,post *models.Post) (*models.Post, error)
	UpdatePost(ctx context.Context,postID int,title,content, image string) (*models.Post, error)
	DeletePost(ctx context.Context,postID int) error
	ListPosts(ctx context.Context) ([]models.Post, error)
	GetUserPosts(ctx context.Context,UserID,limit int64,cursor string)(*dto.PaginatedResponse, error)
	CreateCacheUser(ctx context.Context,userID int ,username,name string)(error)
	GenerateUploadURL(ctx context.Context, filename, contentType string, userID int) (string, map[string]string, error)
	IncreaseLikeCount(ctx context.Context,postID int) error
	DecreaseLikeCount(ctx context.Context,postID int) error

	
}

package repository

import (
	"context"
	"database/sql"

	"github.com/abelmalu/golang-posts/post/internal/models"
)


type PostRepository struct {

	DB *sql.DB
}

func NewPostRepository(DB *sql.DB)*PostRepository{

	return &PostRepository{
		DB:DB, 
	}
	
}

func (pr *PostRepository) CreatePost(ctx context.Context, post *models.Post)(*models.Post,error){

	panic("")
}
func (pr *PostRepository) UpdatePost(ctx context.Context, ID int)(*models.Post,error){

	panic("")
}
func (pr *PostRepository) DeletePost(postID string)(error){

	panic("")
}
func (pr *PostRepository) ListPosts(ctx context.Context)([]models.Post,error){

	panic("")
}

   
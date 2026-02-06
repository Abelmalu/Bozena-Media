package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"github.com/abelmalu/golang-posts/post/internal/models"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(DB *sql.DB) *PostRepository {

	return &PostRepository{
		DB: DB,
	}

}

func (PostRepository *PostRepository) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	
	query := `INSERT INTO posts (title,content,user_id) VALUES($1,$2,$3) RETURNING id`

	err := PostRepository.DB.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID).Scan(
		&post.Id,
	)
	if err != nil {

		log.Printf("Error while inserting a post %v", err)
		return nil, errors.New("Failed to create a post")
	}

	return post, nil
}
func (pr *PostRepository) UpdatePost(ctx context.Context, ID int) (*models.Post, error) {

	panic("")
}
func (pr *PostRepository) DeletePost(postID string) error {

	panic("")
}
func (PostRepository *PostRepository) ListPosts(ctx context.Context) ([]models.Post, error) {

	var posts []models.Post
	query := `SELECT * FROM posts`

	rows, err := PostRepository.DB.Query(query)
	if err != nil {

		log.Printf("error during db query %v", err)
		
		return nil,err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post

		rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserID)
		posts = append(posts, post)

	}

	if err = rows.Err(); err != nil {
		return nil,err
	}

	return posts,nil
}

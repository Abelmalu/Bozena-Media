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
		&post.ID,
	)
	if err != nil {

		log.Printf("Error while inserting a post %v", err)
		return nil, errors.New("Failed to create a post")
	}

	return post, nil
}
func (pr *PostRepository) UpdatePost(ctx context.Context, ID int, title string, content string) (*models.Post, error) {

	var post models.Post

	query := `UPDATE posts SET title=$1, content=$2 WHERE id=$3  RETURNING id, title, content`

	err := pr.DB.QueryRowContext(ctx,query, title, content, ID).
		Scan(&post.ID, &post.Title, &post.Content)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	return &post, nil
}
func (pr *PostRepository) DeletePost(ctx context.Context,postID int)(error) {

	query := `DELETE FROM posts WHERE id=$1`

	result,err := pr.DB.ExecContext(ctx,query,postID)
	
	if err != nil{
		log.Printf("DB erro %v", err)
		
		return nil
 
	}
	rowsAffected,err := result.RowsAffected()
	if err != nil{

		log.Printf("DB erro %v", err)
		
		return nil

	}

	if rowsAffected == 0{
		log.Printf("DB erro rows affected came zero")
	
		return err


	}
	return nil
}
func (PostRepository *PostRepository) ListPosts(ctx context.Context) ([]models.Post, error) {

	var posts []models.Post
	query := `SELECT * FROM posts`

	rows, err := PostRepository.DB.Query(query)
	if err != nil {

		log.Printf("error during db query %v", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post

		rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
		posts = append(posts, post)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

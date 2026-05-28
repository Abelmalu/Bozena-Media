package repository

import (
	"context"
	"database/sql"
	"errors"

	ierrors "github.com/abelmalu/golang-posts/post/internal/errors"
	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
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
	var appErr *pgconn.PgError
	if err != nil {

		if errors.As(err, &appErr) {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}
		if errors.Is(err, context.Canceled) {

			return nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
		}

		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	return post, nil
}
func (pr *PostRepository) UpdatePost(ctx context.Context, ID int, title string, content string) (*models.Post, error) {

	var post models.Post

	query := `UPDATE posts SET title=$1, content=$2 WHERE id=$3  RETURNING id, title, content`

	err := pr.DB.QueryRowContext(ctx, query, title, content, ID).
		Scan(&post.ID, &post.Title, &post.Content)

	var pgErr *pgconn.PgError
	if err != nil {

		if errors.As(err, &pgErr) {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}
		if errors.Is(err, context.Canceled) {

			return nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
		}

		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	return &post, nil
}
func (pr *PostRepository) DeletePost(ctx context.Context, postID int) error {

	query := `DELETE FROM posts WHERE id=$1`

	result, err := pr.DB.ExecContext(ctx, query, postID)

	var appErr *pgconn.PgError
	if err != nil {

		if errors.As(err, &appErr) {

			return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}
		if errors.Is(err, context.Canceled) {

			return ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
		}

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	if rowsAffected == 0 {
		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}
	return nil
}
func (PostRepository *PostRepository) ListPosts(ctx context.Context) ([]models.Post, error) {

	var posts []models.Post
	query := `SELECT * FROM posts`

	rows, err := PostRepository.DB.Query(query)

	var pgErr *pgconn.PgError
	if err != nil {

		if errors.As(err, &pgErr) {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}
		if errors.Is(err, context.Canceled) {

			return nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
		}

		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	defer rows.Close()

	for rows.Next() {
		var post models.Post

		rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
		posts = append(posts, post)

	}

	if err = rows.Err(); err != nil {
		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	return posts, nil
}


func (postRepo *PostRepository) GetUserPosts(ctx context.Context,UserID int64)([]models.Post, error){

	var posts []models.Post

	query :=`SELECT * FROM posts WHERE id=$1`

	rows,err := postRepo.DB.QueryContext(ctx,query,UserID)
	var pgErr *pgconn.PgError
	if err != nil {

		if errors.As(err, &pgErr) {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}
		if errors.Is(err, context.Canceled) {

			return nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
		}

		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	for rows.Next(){
		var post models.Post
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
		posts = append(posts,post)

	} 
	if err = rows.Err(); err != nil {
		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}
	defer rows.Close()
	return nil,nil
}
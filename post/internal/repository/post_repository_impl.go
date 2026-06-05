package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/abelmalu/golang-posts/post/internal/dto"
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

	query := `INSERT INTO posts (title,conte
	nt,user_id) VALUES($1,$2,$3) RETURNING id`

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

	query := ` DELETE FROM posts WHERE id=$1 `

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
	query := `
SELECT id, title, content, user_id
FROM posts
`

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
	defer rows.Close()

	return posts, nil
}

func (postRepo *PostRepository) GetUserPosts(ctx context.Context, UserID int64, limit int64, cursor string) (*dto.PaginatedResponse, error) {

	var posts []*models.Post
	var hasNext bool
	var after string

	if cursor != "" {

		cursorByte, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}

		cursorStr := string(cursorByte)
		cursorInt, err := strconv.Atoi(cursorStr)

		if err != nil {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}

		query := `SELECT id,title,content,user_id FROM posts WHERE user_id=$1 AND id < $2 ORDER BY id DESC LIMIT $3`

		rows, err := postRepo.DB.QueryContext(ctx, query, UserID, cursorInt, (limit + 1))

		if err != nil {
			var pgErr *pgconn.PgError

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

			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

			}

			if len(posts) == int(limit) {

				hasNext = true

				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(post.ID)))
				break

			}

			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}
	} else {

		query := `SELECT id,title,content,user_id FROM posts WHERE user_id=$1 ORDER BY id DESC LIMIT $2`

		rows, err := postRepo.DB.QueryContext(ctx, query, UserID, (limit + 1))

		if err != nil {
			var pgErr *pgconn.PgError

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

			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

			}

			if len(posts) == int(limit) {

				hasNext = true

				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(post.ID)))
				break

			}

			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}

	}

	return &dto.PaginatedResponse{
		Posts:   posts,
		HasNext: hasNext,
		Cursor:  after,
	}, nil
}

func (postRepo *PostRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	query := `INSERT INTO users_cache (user_id,username,name)  VALUES($1,$2,$3)`

	_, err := postRepo.DB.ExecContext(ctx, query, userID, username,name)

	var pgErr *pgconn.PgError
	if err != nil {

		if errors.As(err, &pgErr) {

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
    
	return nil
}

package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"

	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/like/internal/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type LikeRespository struct {
	DB *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRespository {

	return &LikeRespository{
		DB: db,
	}
}

func (likeRespository *LikeRespository) ToggleLike(ctx context.Context, state bool, userID, postID int) (string,error) {

	if state {
		query := `
			INSERT INTO likes (user_id, post_id) 
			VALUES ($1, $2) 
			ON CONFLICT (user_id, post_id) DO NOTHING; `

		_, err := likeRespository.DB.ExecContext(ctx, query, userID, postID)

		var pgErr *pgconn.PgError
		if err != nil {

			if errors.As(err, &pgErr) {

				return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}
			if errors.Is(err, context.Canceled) {

				return "",ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
			}
			if errors.Is(err, context.DeadlineExceeded) {

				return "", ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
			}

			return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

		}


		return "post liked successfully",nil
	} else {

		query := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2;`

		result, err := likeRespository.DB.ExecContext(ctx, query,userID, postID)

		var appErr *pgconn.PgError
		if err != nil {

			if errors.As(err, &appErr) {

				return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}
			if errors.Is(err, context.Canceled) {

				return "",ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
			}
			if errors.Is(err, context.DeadlineExceeded) {

				return "",ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
			}

			return "",ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

		}
		rowsAffected, err := result.RowsAffected()

		if err != nil {
			return "",ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

		}

		if rowsAffected == 0 {
			return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

		}

		return "post unliked successfully",nil


	}


}


func (likeRepository *LikeRespository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	query := `INSERT INTO users_cache (user_id,username,name)  VALUES($1,$2,$3)`

	_, err := likeRepository.DB.ExecContext(ctx, query, userID, username,name)

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



func (likeRepository *LikeRespository) CreateCachePost(ctx context.Context, postID int, title string)error{



	query := `INSERT INTO posts_cache (post_id,title)  VALUES($1,$2)`

	_, err := likeRepository.DB.ExecContext(ctx, query, postID, title)

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


func (likeRepository *LikeRespository) GetPostLikes (ctx context.Context, postID,limit int,cursor string)(*dto.PaginatedPostLikesResponse,error) {

	var usersLiked []*dto.User
	var after string
	var hasNext bool


	if cursor != "" {

		cursorByte,err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}

		cursorStr := string(cursorByte)

		cursorInt,err := strconv.Atoi(cursorStr)

		if err != nil {

			return nil,ierrors.NewDatabaseError(ierrors.MSGDatabaseError,err)
		}

			// Joining likes table and users to get who liked a post 
		query := `SELECT u.user_id,u.name,u.username FROM users_cache u 
				  INNER JOIN likes l ON u.user_id = l.user_id WHERE l.post_id = $1 AND l.id < $2 ORDER BY l.id DESC LIMIT $3 `

		rows, err := likeRepository.DB.QueryContext(ctx,query,postID,cursorInt,(limit + 1))		  

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


		for rows.Next() {
			var userLiked dto.User

			if err := rows.Scan(&userLiked.ID,&userLiked.Name,&userLiked.Username); err != nil {


				return nil,ierrors.NewDatabaseError(ierrors.MSGDatabaseError,err)
			}

			if len(usersLiked) == limit {

				hasNext = true 

				cursor = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(userLiked.ID)))

				break
			}

			usersLiked = append(usersLiked, &userLiked)

		}


		return &dto.PaginatedPostLikesResponse{
			UsersLiked: usersLiked,
			HasNext: hasNext,
			Cursor: after,
		},nil
	} else {

		query := `SELECT u.user_id,u.name,u.username FROM users_cache u 
				  INNER JOIN likes l ON u.user_id = l.user_id WHERE l.post_id = $1  ORDER BY l.id DESC LIMIT $2 `


		rows, err := likeRepository.DB.QueryContext(ctx,query,postID,(limit + 1))		  

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


		for rows.Next() {
			var userLiked dto.User

			if err := rows.Scan(&userLiked.ID,&userLiked.Name,&userLiked.Username); err != nil {


				return nil,ierrors.NewDatabaseError(ierrors.MSGDatabaseError,err)
			}

			if len(usersLiked) == limit {

				hasNext = true 

				cursor = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(userLiked.ID)))

				break
			}

			usersLiked = append(usersLiked, &userLiked)

		}


		return &dto.PaginatedPostLikesResponse{
			UsersLiked: usersLiked,
			HasNext: hasNext,
			Cursor: after,
		},nil
	}
}
 
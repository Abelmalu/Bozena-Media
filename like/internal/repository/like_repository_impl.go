package repository

import (
	"context"
	"database/sql"
	"errors"

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



	query := `INSERT INTO posts_cache (post_id,tile)  VALUES($1,$2)`

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
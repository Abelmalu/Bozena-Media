package repository

import (
	"context"
	"database/sql"
	"errors"

	ierrors "github.com/abelmalu/golang-posts/follow/internal/errors"
	"github.com/jackc/pgx/v5/pgconn"
)


type FollowRepository struct {

	DB *sql.DB
}


func NewFollowRepository(DB *sql.DB)  *FollowRepository {


	return &FollowRepository{
		DB: DB,
	}
}

func (followRepository *FollowRepository) ToggleFollow(ctx context.Context,follow bool,followerID,followingID int)(string,error){


	if follow {
		query := `
			INSERT INTO follows (follower_id, following_id) 
			VALUES ($1, $2) 
			ON CONFLICT (follower_id, following_id) DO NOTHING; `

		_, err := followRepository.DB.ExecContext(ctx, query, followerID, followingID)

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


		return "followed successfully",nil
	} else {

		query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2;`

		result, err := followRepository.DB.ExecContext(ctx, query,followerID, followingID)

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

		return "unfollowed successfully",nil


	}

	
}

package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"

	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/follow/internal/errors"
	"github.com/abelmalu/golang-posts/follow/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

type FollowRepository struct {
	DB *sql.DB
}

func NewFollowRepository(DB *sql.DB) *FollowRepository {

	return &FollowRepository{
		DB: DB,
	}
}

func (followRepository *FollowRepository) ToggleFollow(ctx context.Context, follow bool, followerID, followingID int) (string, error) {

	if follow {
		query := `
			INSERT INTO follows (follower_id, following_id) 
			VALUES ($1, $2); `

		_, err := followRepository.DB.ExecContext(ctx, query, followerID, followingID)

		var pgErr *pgconn.PgError
		if err != nil {

			if errors.As(err, &pgErr) {

				switch pgErr.Code{

				case "23505":
					return "",ierrors.NewBadRequestError(ierrors.ErrorMessage("already following the user"),nil)
				}

				return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}
			if errors.Is(err, context.Canceled) {

				return "", ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
			}
			if errors.Is(err, context.DeadlineExceeded) {

				return "", ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
			}

			return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

		}

		return "followed successfully", nil
	} else {

		query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2;`

		result, err := followRepository.DB.ExecContext(ctx, query, followerID, followingID)

		var appErr *pgconn.PgError
		if err != nil {

			if errors.As(err, &appErr) {

				return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}
			if errors.Is(err, context.Canceled) {

				return "", ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
			}
			if errors.Is(err, context.DeadlineExceeded) {

				return "", ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
			}

			return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

		}
		rowsAffected, err := result.RowsAffected()

		if err != nil {
			return "", ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

		}

		if rowsAffected == 0 {
			return "", ierrors.NewBadRequestError(ierrors.ErrorMessage("Already unfollowed the user"), err)

		}

		return "unfollowed successfully", nil

	}

}

func (followRepository *FollowRepository) GetUserFollowers(ctx context.Context, followingID, limit int, cursor string) (*dto.PaginatedFollowersResponse, error) {

	var followers []*models.User
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

		// query := `SELECT id,follower_id,following_id FROM follows WHERE following_id = $1 AND id < $2 ORDER BY id DESC LIMIT $3`
		query := `
		 SELECT users_cache.user_id, users_cache.name, users_cache.username FROM users_cache 
		 INNER JOIN follows ON users_cache.user_id = follows.follower_id 
		 WHERE follows.following_id = $1 AND users_cache.user_id < $2
		 ORDER BY users_cache.user_id DESC LIMIT $3  `
		 
		rows, err := followRepository.DB.QueryContext(ctx, query, followingID, cursorInt, (limit + 1))

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
			var follower models.User

			if err := rows.Scan(&follower.ID, &follower.Name, &follower.Username); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}

			if len(followers) == limit {

				hasNext = true

				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(follower.ID)))
				break
			}

			followers = append(followers, &follower)
		}

		if err := rows.Err(); err != nil {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}

		return &dto.PaginatedFollowersResponse{
			Followers: followers,
			HasNext:   hasNext,
			Cursor:    after,
		}, nil

	} else {

		//query := ` SELECT id,follower_id,following_id FROM follows WHERE following_id=$1 ORDER BY id DESC LIMIT $2`
		query := `SELECT users_cache.user_id,users_cache.name,users_cache.username from users_cache 
		INNER JOIN follows ON users_cache.user_id = follows.follower_id WHERE follows.following_id=$1 ORDER BY users_cache.user_id DESC LIMIT $2 `

		rows, err := followRepository.DB.QueryContext(ctx, query, followingID, (limit + 1))

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

			var follower models.User

			if err = rows.Scan(&follower.ID, &follower.Name, &follower.Username); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

			}

			if len(followers) == limit {

				hasNext = true

				//changing the id of the follower table to string
				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(follower.ID)))
				break

			}

			followers = append(followers, &follower)

		}

		if err := rows.Err(); err != nil {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}

	}

	return &dto.PaginatedFollowersResponse{
		Followers: followers,
		Cursor:    after,
		HasNext:   hasNext,
	}, nil
}

func (followRepository *FollowRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	query := `INSERT INTO users_cache (user_id,username,name)  VALUES($1,$2,$3)`

	_, err := followRepository.DB.ExecContext(ctx, query, userID, username, name)

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

func (followRepository *FollowRepository) GetUserUserFollowings(ctx context.Context, followerId, limit int, cursor string) (*dto.PaginatedFollowingsResponse, error) {

	var followings []*models.User
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

		query := ` SELECT users_cache.user_id, users_cache.name,users_cache.username FROM users_cache
					INNER JOIN follows ON users_cache.id = following_id 
					WHERE follower_id=$1 AND follows.id < $2 ORDER BY follows.id DESC LIMIT $3
				   
		
		`

		rows, err := followRepository.DB.QueryContext(ctx, query, followerId, cursorInt, (limit + 1))

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

			var following models.User

			if err := rows.Scan(&following.ID, &following.Name, &following.Username); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}

			if len(followings) == limit {

				hasNext = true
				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(following.ID)))
				break

			}

			followings = append(followings, &following)
		}

		return &dto.PaginatedFollowingsResponse{
			Followings: followings,
			HasNext:    hasNext,
			Cursor:     after,
		}, nil

	} else {

		query := ` SELECT u.user_id,u.name,u.username FROM users_cache u
				   INNER JOIN follows f ON u.user_id = f.following_id WHERE f.follower_id=$1 ORDER BY f.id DESC LIMIT $2 `

		rows, err := followRepository.DB.QueryContext(ctx, query, followerId,(limit + 1))

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

			var following models.User

			if err := rows.Scan(&following.ID, &following.Name, &following.Username); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}

			if len(followings) == limit {

				hasNext = true
				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(following.ID)))
				break

			}

			followings = append(followings, &following)

		}

	}

	return &dto.PaginatedFollowingsResponse{
		Followings: followings,
		HasNext:    hasNext,
		Cursor:     after,
	}, nil

}

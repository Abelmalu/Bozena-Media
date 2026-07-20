package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/abelmalu/golang-posts/Feed/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/Feed/internal/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type FeedRepository struct {
	DB *sql.DB
}

func NewFeedRepository(db *sql.DB) *FeedRepository {

	return &FeedRepository{
		DB: db,
	}
}

func (feedRepo *FeedRepository) GetUserFeed(ctx context.Context, cursor string, userID, limit int) (*dto.PaginatedResponse, error) {

	var userFeeds []*dto.UserFeed
	var after string
	var hasNext bool

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

		query := `

				 SELECT fe.id,fe.owner_id,u.username,u.name,p.post_id,p.title,p.content,p.like_count,p.image FROM feed_entries fe
				 INNER JOIN  users_cache u ON  fe.owner_id = u.user_id 
				 INNER JOIN  posts_cache p ON   fe.post_id = p.post_id 
				 WHERE fe.user_id = $1 AND fe.id < $2 ORDER BY fe.id DESC LIMIT $3

		`

		rows, err := feedRepo.DB.QueryContext(ctx, query, userID, cursorInt, (limit + 1))

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

			var userFeed dto.UserFeed

			if err := rows.Scan(&userFeed.ID, &userFeed.PostOwnerID, &userFeed.UserName, &userFeed.Name, &userFeed.PostID,&userFeed.PostTitle, &userFeed.PostContent,&userFeed.LikeCount,&userFeed.Image); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}

			if limit == len(userFeeds) {

				hasNext = true

				cursor = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(userFeed.ID)))

				break

			}

			userFeeds = append(userFeeds, &userFeed)

		}

		return &dto.PaginatedResponse{

			UserFeeds: userFeeds,
			Cursor:    after,
			HasNext:   hasNext,
		}, nil

	} else {

		query := `
				 SELECT fe.id,fe.owner_id,u.username,u.name,p.post_id,p.title,p.content,p.like_count,p.image FROM feed_entries fe
				 INNER JOIN  users_cache u ON  fe.owner_id = u.user_id 
				 INNER JOIN  posts_cache p ON   fe.post_id = p.post_id 
				 WHERE fe.user_id = $1 ORDER BY fe.id DESC LIMIT $2
		
		
		`

		rows, err := feedRepo.DB.QueryContext(ctx, query, userID, (limit + 1))

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

			var userFeed dto.UserFeed

			if err := rows.Scan(&userFeed.ID, &userFeed.PostOwnerID, &userFeed.UserName, &userFeed.Name, &userFeed.PostID,&userFeed.PostTitle, &userFeed.PostContent,&userFeed.LikeCount,&userFeed.Image); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}

			if limit == len(userFeeds) {

				hasNext = true

				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(userFeed.ID)))

				break

			}

			userFeeds = append(userFeeds, &userFeed)

		}

		return &dto.PaginatedResponse{
			UserFeeds: userFeeds,
			Cursor:    after,
			HasNext:   hasNext,
		}, nil
	}

}

func (feedRepository *FeedRepository) CreateCachePost(ctx context.Context, postID int, title, content string) error {

	query := `INSERT INTO posts_cache (post_id,title,content)  VALUES($1,$2,$3)`

	_, err := feedRepository.DB.ExecContext(ctx, query, postID, title, content)

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

func (feedRepository *FeedRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	query := `INSERT INTO users_cache (user_id,username,name)  VALUES($1,$2,$3)`

	_, err := feedRepository.DB.ExecContext(ctx, query, userID, username, name)

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
func (repo *FeedRepository) InsertFeedEntries(ctx context.Context, followersID []int, postID, ownerID int) error {
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO feed_entries(user_id, post_id, owner_id) 
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, followerID := range followersID {
		_, err := stmt.ExecContext(ctx, followerID, postID, ownerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
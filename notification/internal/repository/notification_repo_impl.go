package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/abelmalu/golang-posts/notification/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/notification/internal/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {

	return &NotificationRepository{
		DB: db,
	}
}

func (notificationRepo *NotificationRepository) GetUserNotifications(ctx context.Context, userID int, cursor string, limit int) (*dto.PaginatedResponse, error) {

	var UserNotifications []*dto.UserNotification
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

		query := ` SELECT n.id,u.username,n.actor_id,n.message,n.created_at FROM notifications AS n
				INNER JOIN users_cache AS u ON n.actor_id = u.user_id WHERE recipient_id = $1 AND n.id < $2 ORDER BY id DESC limit $3 `

		rows, err := notificationRepo.DB.QueryContext(ctx, query, userID, cursorInt, (limit + 1))

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

			var UserNotification dto.UserNotification

			if err := rows.Scan(&UserNotification.ID, &UserNotification.UseraName, &UserNotification.ActorID, &UserNotification.Message, &UserNotification.CreatedAT); err != nil {
				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

			}

			if len(UserNotifications) == limit {

				hasNext = true
				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(UserNotification.ID)))

				break

			}

			UserNotifications = append(UserNotifications, &UserNotification)

		}

		return &dto.PaginatedResponse{

			HasNext:           hasNext,
			Cursor:            after,
			UserNotifications: UserNotifications,
		}, nil
	} else {

		query := ` SELECT n.id,u.username,n.actor_id,n.message,n.created_at FROM notifications AS n
				INNER JOIN users_cache AS u ON n.actor_id = u.user_id WHERE recipient_id = $1  ORDER BY id DESC limit $2 `

		rows, err := notificationRepo.DB.QueryContext(ctx, query, userID, (limit + 1))

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

			var UserNotification dto.UserNotification

			if err := rows.Scan(&UserNotification.ID, &UserNotification.UseraName, &UserNotification.ActorID, &UserNotification.Message, &UserNotification.CreatedAT); err != nil {
				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

			}

			if len(UserNotifications) == limit {

				hasNext = true
				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(UserNotification.ID)))

				break

			}

			UserNotifications = append(UserNotifications, &UserNotification)

		}

		return &dto.PaginatedResponse{

			HasNext:           hasNext,
			Cursor:            after,
			UserNotifications: UserNotifications,
		}, nil

	}

}

func (notificationRepo *NotificationRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {

	query := `INSERT INTO users_cache (user_id,username,name)  VALUES($1,$2,$3)`

	_, err := notificationRepo.DB.ExecContext(ctx, query, userID, username, name)

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



func (notificationRepo *NotificationRepository)	InsertUserNotification(ctx context.Context,actorID,recipientID int)  error {


query := `INSERT INTO notifications (actor_id,recipient_id)  VALUES($1,$2)`

	_, err := notificationRepo.DB.ExecContext(ctx, query, actorID, recipientID)

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

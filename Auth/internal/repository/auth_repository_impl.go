package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/abelmalu/golang-posts/Auth/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/Auth/internal/errors"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
	"github.com/abelmalu/golang-posts/Auth/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository struct {
	DB     *sql.DB
	logger *platform.Logger
}

func NewAuthRepository(db *sql.DB) *AuthRepository {

	return &AuthRepository{
		DB: db,
	}
}

func (authRepo *AuthRepository) Register(ctx context.Context, user *model.User) (*model.User, error) {
	var newUser model.User

	query := `INSERT INTO users(name,username,email,password) VALUES($1,$2,$3,$4) RETURNING id,role,name,username`
	if err := authRepo.DB.QueryRowContext(ctx, query, user.Name, user.Username, user.Email, user.Password).Scan(&newUser.ID, &newUser.Role, &newUser.Name, &newUser.Username); err != nil {
		// Change *pq.Error to *pgconn.PgError
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {

			switch pgErr.Code {

			case "23505":
				switch pgErr.ConstraintName {
				case "users_username_key":
					return nil, ierrors.NewValidationError(ierrors.MSGUsenameAlreadyExists, nil, err)
				case "users_email_key":
					return nil, ierrors.NewValidationError(ierrors.MSGEmailAlreadyExists, nil, err)

				default:
					return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

				}
			default:
				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
			}

		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)

		}
		if errors.Is(err, context.Canceled) {
			return nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)

		}

		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	return &newUser, nil
}
func (authrepo *AuthRepository) Login(ctx context.Context, userName, password string) (*model.User, error) {

	var user model.User
	query := `SELECT * FROM users WHERE username=$1`
	if err := authrepo.DB.QueryRowContext(ctx, query, userName).Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role, &user.FailedLoginAttempts, &user.IsPermanentlyLocked, &user.TemporaryLockUntil, &user.FollowerCount, &user.FollowingCount,&user.Avatar); err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

		}

		if errors.Is(err, sql.ErrNoRows) {

			return nil, ierrors.NewNotFoundError(ierrors.MSGUserNotFound, err)

		}

		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	return &user, nil

}
func (authRepo *AuthRepository) Logout(ctx context.Context, tokenID string) error {

	query := `DELETE FROM refresh_tokens WHERE token_text=$1`

	result, err := authRepo.DB.ExecContext(ctx, query, tokenID)

	var pgErr *pgconn.PgError
	if err != nil {

		if errors.As(err, &pgErr) {

			return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}

		if errors.Is(err, context.DeadlineExceeded) {

			return ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)

		}
		if errors.Is(err, context.Canceled) {

			return ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {

			return ierrors.NewNotFoundError(ierrors.MSGNotFound, err)

		}

	}

	if rowsAffected == 0 {

		return ierrors.NewNotFoundError(ierrors.MSGNotFound, err)

	}

	return err
}
func (authRepo *AuthRepository) StoreRefreshTokens(userID int, refreshToken string, expiresAt time.Time, clientType string) (sql.Result, error) {

	// hashing the token before inserting to a db
	refreshToken = utils.HashToken(refreshToken)

	query := `INSERT INTO refresh_tokens (user_id,token_text,expires_at,client_type) VALUES($1,$2,$3,$4)`

	result, err := authRepo.DB.Exec(query, userID, refreshToken, expiresAt, clientType)
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

	return result, nil

}

func (authRepo *AuthRepository) RevokeRefreshToken(refreshToken string) error {

	query := `
	
	UPDATE refresh_tokens SET revoked=TRUE 
	WHERE token_text = $1 AND revoked=FALSE `

	result, err := authRepo.DB.Exec(query, refreshToken)

	rowsAffected, err := result.RowsAffected()
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

	// detect reuse attempt
	if rowsAffected == 0 {
		// token was already revoked or doesn't exist
		if errors.Is(err, sql.ErrNoRows) {

			return ierrors.NewNotFoundError(ierrors.MSGNotFound, err)
		}

	}

	return nil

}

func (authRepo *AuthRepository) GetRefreshToken(refreshToken string) (*model.RefreshToken, error) {

	var refreshRecord model.RefreshToken

	// hashing the token because stored tokens are hashed
	hashedrefreshToken := utils.HashToken(refreshToken)

	query := `SELECT * FROM refresh_tokens where token_text = $1;`

	if err := authRepo.DB.QueryRow(query, hashedrefreshToken).Scan(&refreshRecord); err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, ierrors.NewNotFoundError(ierrors.MSGNotFound, err)
		}
		if errors.Is(err, context.Canceled) {

			return nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
		}
		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	return &refreshRecord, nil
}

func (authRepo *AuthRepository) GetUserByID(ID int) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id=$1`

	err := authRepo.DB.QueryRow(query, ID).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, ierrors.NewNotFoundError(ierrors.MSGNotFound, err)
		}

		if errors.Is(err, context.Canceled) {

			return nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
		}
		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}
	return &user, nil
}

func (authRepo *AuthRepository) SearchUser(ctx context.Context, username, cursor string, limit int) (*dto.PaginatedResponse, error) {

	var users []*model.User
	var after string
	var hasNext bool

	if cursor != "" {

		cursorByte := base64.StdEncoding.EncodeToString([]byte(cursor))
		cursorStr := string(cursorByte)
		cursorInt, err := strconv.Atoi(cursorStr)
		username = "% "
		if err != nil {

			return nil, ierrors.NewValidationError(ierrors.MSGBadRequest, nil, err)
		}

		query := ` SELECT id,name,username FROM users WHERE username ILIKE $1 AND id < $2 ORDER BY id DESC LIMIT $3 `

		rows, err := authRepo.DB.QueryContext(ctx, query, username, cursorInt, (limit + 1))

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
			var user model.User

			if err := rows.Scan(&user.ID, &user.Name, &user.Username); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

			}

			if len(users) == limit {

				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(user.ID)))
				hasNext = true

				break

			}

			users = append(users, &user)

		}

	} else {

		query := ` SELECT id,name,username FROM users WHERE username ILIKE $1 ORDER BY id DESC LIMIT $2`

		rows, err := authRepo.DB.QueryContext(ctx, query, username, (limit + 1))

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

			var user model.User

			if err = rows.Scan(&user.ID, &user.Name, &user.Username); err != nil {

				return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

			}

			if len(users) == limit {

				hasNext = true

				//changing the id of the follows table to string
				after = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(user.ID)))
				break

			}

			users = append(users, &user)

		}

		if err := rows.Err(); err != nil {

			return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
		}

	}

	return &dto.PaginatedResponse{
		Users:   users,
		Cursor:  after,
		HasNext: hasNext,
	}, nil
}

func (authRepo *AuthRepository) UpdateFailedLoginAttempts(ctx context.Context, user *model.User) (*model.User, error) {

	query := ` UPDATE users SET failed_attempts=failed_attempts+1 WHERE id=$1 RETURNING failed_attempts`

	err := authRepo.DB.QueryRowContext(ctx, query, user.ID).Scan(&user.FailedLoginAttempts)

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

	return user, nil
}

func (authRepo *AuthRepository) TemporaryLockUntil(ctx context.Context, user *model.User) (*model.User, error) {

	query := ` UPDATE users SET failed_attempts = $1, temporary_locked_until = $2  WHERE id=$3 RETURNING failed_attempts`

	err := authRepo.DB.QueryRowContext(ctx, query, user.FailedLoginAttempts, user.TemporaryLockUntil, user.ID).Scan(&user.FailedLoginAttempts)

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

	return user, nil

}

func (authRepo *AuthRepository) IncreaseFollowCount(ctx context.Context, followerID, followingID int) error {

	tx, err := authRepo.DB.BeginTx(ctx, nil)

	if err != nil {

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	defer tx.Rollback()

	query1 := `UPDATE users SET follower_count=follower_count+1 WHERE id=$1`

	_, err = tx.ExecContext(ctx, query1, followingID)

	if err != nil {

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	query2 := `UPDATE users SET following_count=following_count+1 WHERE id=$1`

	_, err = tx.ExecContext(ctx, query2, followerID)

	if err != nil {

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	if err := tx.Commit(); err != nil {

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	return nil
}

func (authRepo *AuthRepository) DecreaseFollowCount(ctx context.Context, followerID, followingID int) error {

	tx, err := authRepo.DB.BeginTx(ctx, nil)

	if err != nil {

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	defer tx.Rollback()

	query1 := `UPDATE users SET follower_count=follower_count-1 WHERE id=$1`

	_, err = tx.ExecContext(ctx, query1, followingID)

	if err != nil {

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	query2 := `UPDATE users SET following_count=following_count-1 WHERE id=$1`

	_, err = tx.ExecContext(ctx, query2, followerID)

	if err != nil {

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	if err := tx.Commit(); err != nil {

		return ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)
	}

	return nil
}

func (authRepo *AuthRepository) GetUserProfile(ctx context.Context, userID int64) (*dto.UserProfileResponse, error) {

	var user dto.UserProfileResponse

	query := ` SELECT id,name,username,avatar FROM users WHERE id=$1`

	if err := authRepo.DB.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Name, &user.UserName, &user.Avatar); err != nil {

		if errors.Is(err, sql.ErrNoRows) {

			return nil, ierrors.NewNotFoundError(ierrors.MSGNotFound, err)
		}

		if errors.Is(err, context.Canceled) {

			return nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)
		}
		return nil, ierrors.NewDatabaseError(ierrors.MSGDatabaseError, err)

	}

	return &user, nil

}

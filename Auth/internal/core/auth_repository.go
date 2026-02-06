package core

import (
	"context"
	"database/sql"
	"time"

	model "github.com/abelmalu/golang-posts/Auth/internal/models"
)

type AuthRepository interface {
	Register(ctx context.Context,user *model.User)(*model.User,error)
    StoreRefreshTokens(userID int, refreshToken string, expiresAt time.Time, clientType string) (sql.Result, error)
}
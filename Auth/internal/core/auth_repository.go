package core

import (
	"context"
	"database/sql"
	"time"

	model "github.com/abelmalu/golang-posts/Auth/internal/models"
)

type AuthRepository interface {
	Register(ctx context.Context,user *model.User)(*model.User,error)
	Login(ctx context.Context,userName,password string)(*model.User,error)
	Logout(ctx context.Context, tokenID string) (error)
    StoreRefreshTokens(userID int, refreshToken string, expiresAt time.Time, clientType string) (sql.Result, error)
	RevokeRefreshToken(refreshToken string) error
	GetRefreshToken(refreshToken string) (*model.RefreshToken, error)
	GetUserByID(ID int) (*model.User, error)
}
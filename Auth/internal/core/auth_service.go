package core

import (
	"context"

	"github.com/abelmalu/golang-posts/Auth/internal/dto"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, post *model.User) (*model.User, *model.TokenPair, error)
	Login(ctx context.Context, userName, password string) (*model.User, *model.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshHandler(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	SearchUser(ctx context.Context, username, cursor string, limit int) (*dto.PaginatedResponse, error)
	IncreaseFollowCounts(ctx context.Context, followerID, followingID int) error
	DecreaseFollowCounts(ctx context.Context, followerID, followingID int) error	
	GetUserProfile(ctx context.Context,userID int64)(*model.User,error)

	// UpdateProfile(ctx context.Context,int)

}

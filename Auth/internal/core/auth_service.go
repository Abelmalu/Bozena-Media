package core

import (
	"context"

	model "github.com/abelmalu/golang-posts/Auth/internal/models"
)

type AuthService interface {
	Register(ctx context.Context,post *model.User) (*model.User,*model.TokenPair, error)
}
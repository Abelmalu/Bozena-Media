package core

import (
	"context"

	model "github.com/abelmalu/golang-posts/Auth/internal/models"
)

type AuthRepository interface {
	Register(ctx context.Context,user model.User)(*model.User,error)
}
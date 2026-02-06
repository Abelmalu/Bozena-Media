package service

import (
	"context"

	"github.com/abelmalu/golang-posts/Auth/internal/core"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
)

type AuthService struct {
	repo core.AuthRepository
}


func NewAuthService(authRepo core.AuthRepository) *AuthService{

	return &AuthService{repo:authRepo}
}
func (authSer *AuthService) Register(ctx context.Context, post *model.User) (*model.User, error){

panic("")

}

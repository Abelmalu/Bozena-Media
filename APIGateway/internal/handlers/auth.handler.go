package handler

import (
	"context"

	"github.com/abelmalu/golang-posts/Auth/proto/pb"
)

type AuthService interface{
	Register(ctx context.Context,userName,name,email,password string)(*pb.RegisterResponse,error)
	Login(ctx context.Context,userName,password string)(*pb.LoginResponse,error)
	Logout(ctx context.Context,userName,password string)(*pb.LogoutResponse,error)
}
type AuthHandler struct {

	client AuthService
}


func NewAtuhHandler(au AuthService)*AuthHandler{

	return &AuthHandler{client:au}
}

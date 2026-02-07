package handler

import (
	"context"
	"errors"
	"log"

	"github.com/abelmalu/golang-posts/Auth/internal/core"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct{

	pb.UnimplementedAuthServiceServer
	service core.AuthService

}

func NewAuthHandler(authServ core.AuthService) *AuthHandler{

	return &AuthHandler{service: authServ}

}
// Register registers a new user into the system 
func (authHandler *AuthHandler) Register(ctx context.Context,req *pb.RegisterRequest)(*pb.RegisterResponse,error){

	user := model.User{
		Name: req.Name,
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,

	}

	createdUser,tokens,err := authHandler.service.Register(ctx,&user)
	if err != nil {

        log.Printf("CreateUser failed: %v", err)

        // Map errors to gRPC status codes
        if errors.Is(err, context.Canceled) {
            return nil, status.Error(codes.Canceled, "request canceled")
        }

        if errors.Is(err, context.DeadlineExceeded) {
            return nil, status.Error(codes.DeadlineExceeded, "timeout")
        }

        

        return nil, status.Error(codes.Internal, "internal server error")
    }

	return &pb.RegisterResponse{
		Name: createdUser.Name,
		Username: createdUser.Username,
		Email: createdUser.Email,
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,

	},nil

	
}

func (authHandler *AuthHandler)Login(ctx context.Context, req  *pb.LoginRequest) (*pb.LoginResponse, error){
	user,tokens, err := authHandler.service.Login(ctx,req.Username,req.Password)
	if err != nil{

		return  nil, status.Error(codes.Canceled, "user already exists") 
	}

	return &pb.LoginResponse{
		Id:int64(user.ID),
		Username: user.Username,
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,

	},nil
}
func (authHandler *AuthHandler)Logout(context.Context, *emptypb.Empty) (*pb.LogoutResponse, error){
	panic("")
}
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
)

type AuthHandler struct{

	pb.UnimplementedAuthServiceServer
	service core.AuthService

}

func NewAuthHandler(authServ core.AuthService) *AuthHandler{

	return &AuthHandler{service: authServ}

}
// Register registers a new user into the system 
func (ah *AuthHandler) Register(ctx context.Context,req *pb.RegisterRequest)(*pb.RegisterResponse,error){

	user := model.User{
		Name: req.Name,
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,

	}

	createdUser,err := ah.service.Register(ctx,&user)
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
		AccessToken: "createdUser.AccessToken",
		RefreshToken: "createdUser.RefreshToken",

	},nil

	
}
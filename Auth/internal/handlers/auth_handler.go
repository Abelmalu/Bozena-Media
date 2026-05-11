package handler

import (
	"context"
	"errors"

	"github.com/abelmalu/golang-posts/Auth/internal/core"
	ierrors "github.com/abelmalu/golang-posts/Auth/internal/errors"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct{

	pb.UnimplementedAuthServiceServer
	service core.AuthService
	logger *platform.Logger


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

	var appErr *ierrors.AppError
	if errors.As(err, &appErr) {
      switch appErr.Type {
      case ierrors.TypeValidation:
          return nil, status.Error(codes.InvalidArgument, string(appErr.Message))
      case ierrors.TypeConflict:
          return nil, status.Error(codes.AlreadyExists, string(appErr.Message))
      case ierrors.TypeUnauthorized:
          return nil, status.Error(codes.Unauthenticated, string(appErr.Message))
      case ierrors.TypeNotFound:
          return nil, status.Error(codes.NotFound, string(appErr.Message))
      case ierrors.TypeTimeout:
          return nil, status.Error(codes.DeadlineExceeded, string(appErr.Message))
      default:
          return nil, status.Error(codes.Internal, "internal error")
      }
  }
	

	return &pb.RegisterResponse{
		Name: createdUser.Name,
		Username: createdUser.Username,
		Email: createdUser.Email,
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,

	},nil

	
}

func (authHandler *AuthHandler) Login(ctx context.Context, req  *pb.LoginRequest) (*pb.LoginResponse, error){
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
func (authHandler *AuthHandler)Logout(ctx context.Context, req *emptypb.Empty) (*pb.LogoutResponse, error){
	var refreshToken string
	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return nil, errors.New("Unknown device type")
	}
	values := md.Get("refreshToken")
	if len(values) > 0 {
		refreshToken = values[0]
	} else {

		return nil, errors.New("Unknown device type")

	}

	authHandler.service.Logout(ctx,refreshToken)

	return &pb.LogoutResponse{
		Message: "Successfully Logged Out!",
	},nil
	
}


func (authHandler *AuthHandler) RefreshHandler(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error){

	// var clientType string
	// md, exists := metadata.FromIncomingContext(ctx)

	// if !exists {
	// 	return nil, errors.New("Unknown device type")
	// }
	// values := md.Get("refreshToken")
	// if len(values) > 0 {
	// 	clientType = values[0]
	// } else {

	// 	return nil, errors.New("Unknown device type")

	// }

	tokens,err := authHandler.service.RefreshHandler(ctx,req.RefreshToken)
	if err != nil {

		return nil,err
	}

	return &pb.RefreshResponse{
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	},nil

}

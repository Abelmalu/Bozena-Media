package client

import (
	"context"
	"time"

	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"github.com/abelmalu/golang-posts/internal/models"
	"google.golang.org/grpc"
     "google.golang.org/protobuf/types/known/emptypb"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {

	return &AuthClient{
		client: pb.NewAuthServiceClient(conn),
	}
}

func (ac *AuthClient) Register(ctx context.Context,user *models.User)(*pb.RegisterResponse,error){

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

return ac.client.Register(ctx,&pb.RegisterRequest{
	Name:user.Username,
	Username: user.Username,
	Email: user.Email,
	Password: user.Password,
})

}


func (ac *AuthClient) Login(ctx context.Context,userName,password string)(*pb.LoginResponse,error){

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

return ac.client.Login(ctx,&pb.LoginRequest{
	Username: userName,
	Password: password,
})

}

func (ac *AuthClient) Logout(ctx context.Context,userName,password string)(*pb.LogoutResponse,error){

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

return ac.client.Logout(ctx,&emptypb.Empty{})

}

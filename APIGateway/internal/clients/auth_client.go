package client

import (
	"context"
	"time"

	"github.com/abelmalu/golang-posts/Auth/proto/pb"
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

func (ac *AuthClient) Register(ctx context.Context,userName,name,email,password string)(*pb.RegisterResponse,error){

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

return ac.client.Register(ctx,&pb.RegisterRequest{
	Name:name,
	Username: userName,
	Email: email,
	Password: password,
})

}


func (ac *AuthClient) Login(ctx context.Context,userName,password string)(*pb.LoginResponse,error){

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

return ac.client.Login(ctx,&pb.LoginRequest{
	Username: userName,
	Password: password,
})

}

func (ac *AuthClient) Logout(ctx context.Context)(*pb.LogoutResponse,error){

ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

return ac.client.Logout(ctx,&emptypb.Empty{})

}

func (ac *AuthClient) RefreshHandler(ctx context.Context,refreshToken string)(*pb.RefreshResponse,error){
    ctx,cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()

	return ac.client.RefreshHandler(ctx,&pb.RefreshRequest{
		RefreshToken: refreshToken,
	})
}

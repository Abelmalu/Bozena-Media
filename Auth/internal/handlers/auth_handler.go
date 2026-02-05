package handler

import (
	"context"

	"github.com/abelmalu/golang-posts/Auth/proto/pb"
)

type AuthHandler struct{

	pb.UnimplementedAuthServiceServer

}
\
// Register registers a new user into the system 
func (ah *AuthHandler) Register(ctx context.Context,req *pb.RegisterRequest)(*pb.RegisterRequest,error){

	
}
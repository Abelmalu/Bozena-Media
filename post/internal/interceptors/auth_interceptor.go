package interceptors

import (
	"context"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {

		log.Printf("inside of the auth interceptor")

        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "missing metadata")
        }

        values := md.Get("user-id")
        if len(values) == 0 {
            return nil, status.Error(codes.Unauthenticated, "user-id not provided")
        }
        
		log.Printf("the userId is %v",values[0])
        userID, err := strconv.Atoi(values[0])
        if err != nil {
            return nil, status.Error(codes.Unauthenticated, "invalid user-id")
        }

      
        ctx = context.WithValue(ctx, "userID", userID)

        return handler(ctx, req)
    }
}
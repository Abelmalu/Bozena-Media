package interceptors

import (
	"context"
	"strconv"

	"github.com/abelmalu/golang-posts/platform"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(logger *platform.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Error("couldn't get metadata from context")
			return nil, status.Error(codes.Internal, "Something went wrong")
		}

		values := md.Get("user-id")
		if len(values) == 0 {
			logger.Error("couldn't get user-id from context")

			return nil, status.Error(codes.Unauthenticated, "you are not authenticated")
		}

		userID, err := strconv.Atoi(values[0])
		if err != nil {
			logger.Error("couldn't parse user-id to integer")

			return nil, status.Error(codes.Unauthenticated, "You are not authenticated")
		}

		ctx = context.WithValue(ctx, "userID", userID)

		return handler(ctx, req)
	}
}

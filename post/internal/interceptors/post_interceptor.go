package interceptors

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PostOwnershipInterceptor(db *sql.DB) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if info.FullMethod == "/post.PostService/UpdatePost" || info.FullMethod == "/post.PostService/DeletePost" {

			userID, ok := ctx.Value("userID").(int)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "user identity not found")
			}


			type postRequest interface{ GetId() int32 }
			pReq, ok := req.(postRequest)
			if !ok {
				return nil, status.Error(codes.Internal, "failed to parse request id")
			}
			postID := pReq.GetId()

			
			var ownerID int
			query := "SELECT user_id FROM posts WHERE id = $1"
			err := db.QueryRowContext(ctx, query, postID).Scan(&ownerID)

			if err == sql.ErrNoRows {
				return nil, status.Error(codes.NotFound, "post not found")
			}
			if err != nil {
				return nil, status.Error(codes.Internal, "database error")
			}

			
			if int32(ownerID) != int32(userID) {
				return nil, status.Error(codes.PermissionDenied, "you do not own this post")
			}
		}

		
		return handler(ctx, req)
	}
}
package interceptors

import (
	"context"
	"database/sql"
	"log"

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
			log.Printf("insnide of  postowner interceptor")
			log.Printf(info.FullMethod)


		if info.FullMethod == "/postservice.PostService/UpdatePost" || info.FullMethod == "/postservice.PostService/DeletePost" {
			log.Printf("in the if close")

			userID, ok := ctx.Value("userID").(int)
			log.Printf("user id from postowner interceptor %v",userID)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "user identity not found")
			}


			type postRequest interface{ GetPostId() int64 }
			pReq, ok := req.(postRequest)
			if !ok {
				log.Printf("error in postRequest")
				return nil, status.Error(codes.Internal, "failed to parse request id")
			}
			postID := pReq.GetPostId()

			
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
				log.Printf("u don't own the post")
				return nil, status.Error(codes.PermissionDenied, "you do not own this post")
			}
		}

			log.Printf("below the if close")


		
		return handler(ctx, req)
	}
}
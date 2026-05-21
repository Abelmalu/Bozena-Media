package interceptors

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/abelmalu/golang-posts/platform"
	ierrors "github.com/abelmalu/golang-posts/post/internal/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PostOwnershipInterceptor(db *sql.DB,logger *platform.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
			

		// apply this middleware on this routes
		if info.FullMethod == "/postservice.PostService/UpdatePost" || info.FullMethod == "/postservice.PostService/DeletePost" {

			userID, ok := ctx.Value("userID").(int)

			if !ok {
			logger.Error("couldn't get userID from context")

				return nil, status.Error(codes.Unauthenticated, "user identity not found")
			}


			type postRequest interface{ GetPostId() int64 }
			pReq, ok := req.(postRequest)
			if !ok {
				logger.Error("type assertion failed o")
				return nil, status.Error(codes.Internal, "failed to parse request id")
			}
			postID := pReq.GetPostId()

			
			var ownerID int

			query := "SELECT user_id FROM posts WHERE id = $1"
			err := db.QueryRowContext(ctx, query, postID).Scan(&ownerID)

			if err == sql.ErrNoRows {
				logger.Error(fmt.Sprintf("couldn't find post with %d",postID),zap.Error(err))
				
			return nil, ierrors.NewNotFoundError(ierrors.MSGNotFound, err)

			}
			if err != nil {
				logger.Error("Database Error",zap.Error(err))

				return nil, status.Error(codes.Internal, "Something went wrong")
			}

			
			if int32(ownerID) != int32(userID) {
				logger.Info("Unauthorized attemt to update/delete a post")
				return nil, status.Error(codes.PermissionDenied, "you do not own this post")
			}
		}


		
		return handler(ctx, req)
	}
}
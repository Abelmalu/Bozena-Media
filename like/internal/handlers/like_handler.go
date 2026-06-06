package handler

import (
	"context"
	"errors"

	"github.com/abelmalu/golang-posts/like/internal/core"
	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/like/internal/errors"
	"github.com/abelmalu/golang-posts/like/pkg/utils"
	"github.com/abelmalu/golang-posts/like/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LikeHandler struct {
	pb.UnimplementedLikeServiceServer

	likeService core.LikeService
	logger      *platform.Logger
}

func NewLikeHandler(likeService core.LikeService, logger *platform.Logger) *LikeHandler {

	return &LikeHandler{

		likeService: likeService,
		logger:      logger,
	}
}

func (likeHandler *LikeHandler) ToggleLike(ctx context.Context, req *pb.LikeRequest) (*pb.LikeResponse, error) {
	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "something went wrong")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "something went wrong")

	}
	userID, err := utils.GetUserID(ctx)
	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "something went wrong")

	}
	if errors.Is(err, ierrors.ErrUSerIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "something went wrong")

	}

	postID, err := utils.GetPostID(ctx)
	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "something went wrong")

	}
	if errors.Is(err, ierrors.ErrPostIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "something went wrong")

	}

	likeRequest := dto.ToggleLikeRequest{
		State: &req.State,
	}
	validationErrStr := likeRequest.Validate()

	if validationErrStr != "" {

		return nil, status.Error(codes.InvalidArgument, validationErrStr)

	}

	var appErr *ierrors.AppError
	likeResp, err := likeHandler.likeService.ToggleLike(ctx, req.State, userID, postID)

	if err != nil {

		likeHandler.logger.Error("Like service error", zap.Error(err), zap.String("requestID", requestID))

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
			case ierrors.TypeCancelled:
				return nil, status.Error(codes.Canceled, string(appErr.Message))

			default:
				return nil, status.Error(codes.Internal, "internal error")
			}

		}

		if errors.Is(err, context.Canceled) {

			return nil, status.Error(codes.Canceled, "Request canceled")
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, status.Error(codes.DeadlineExceeded, "Request timeout")

		}

	}
	return &pb.LikeResponse{

		Message: likeResp.Message,
	}, nil

}

func (likeHandler *LikeHandler)	GetPostLikes(ctx context.Context, req *pb.GetPostLikesRequest) (*pb.GetPostLikesResponse, error) {

	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "something went wrong")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "something went wrong")

	}
	resp,err := likeHandler.likeService.GetPostLikes(ctx,int(req.PostId),int(req.Limit),req.Cursor)


	if err != nil {

		likeHandler.logger.Error("Like service error", zap.Error(err), zap.String("requestID", requestID))
		
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
			case ierrors.TypeCancelled:
				return nil, status.Error(codes.Canceled, string(appErr.Message))

			default:
				return nil, status.Error(codes.Internal, "internal error")
			}

		}

		if errors.Is(err, context.Canceled) {

			return nil, status.Error(codes.Canceled, "Request canceled")
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, status.Error(codes.DeadlineExceeded, "Request timeout")

		}

	}

	pbUsersLiked := make([]*pb.User,0,len(resp.UsersLiked))


	for _,userLiked := range resp.UsersLiked {


		pbUser := &pb.User{
			Name: userLiked.Name,
			Username: userLiked.Username,
		}

		pbUsersLiked = append(pbUsersLiked, pbUser)
	}

	return &pb.GetPostLikesResponse{
		Users: pbUsersLiked,
		HasNext: resp.HasNext,
		Cursor: resp.Cursor,

	},nil
	
}


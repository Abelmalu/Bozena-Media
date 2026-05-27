package handler

import (
	"context"
	"errors"

	"github.com/abelmalu/golang-posts/follow/internal/core"
	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/follow/internal/errors"
	"github.com/abelmalu/golang-posts/follow/pkg/utils"
	"github.com/abelmalu/golang-posts/follow/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FollowHandler struct {
	followService core.FollowService
	logger        *platform.Logger
	pb.UnimplementedFollowServiceServer
}

func NewFollowHandler(followService core.FollowService, logger *platform.Logger) *FollowHandler {

	return &FollowHandler{
		followService: followService,
		logger:        logger,
	}
}

func (followHandler *FollowHandler) ToggleFollow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowResponse, error) {
	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {
		followHandler.logger.Error("meta data not found",zap.Error(err),zap.String("requestID",requestID))
		return nil, status.Error(codes.Internal, "something went wrong")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {
		followHandler.logger.Error("requestID not found",zap.Error(err),zap.String("requestID",requestID))

		return nil, status.Error(codes.InvalidArgument, "something went wrong")

	}
	var followRequest = dto.FollowRequest{

		Follow: &req.Follow,
	}

	validationErrStr := followRequest.Validate()

	if validationErrStr != "" {

		followHandler.logger.Error("Follow service", zap.Error(err), zap.String("requestID", requestID))

		return nil, status.Error(codes.InvalidArgument, validationErrStr)

	}

	resp, err := followHandler.followService.ToggleFollow(ctx, *followRequest.Follow, int(req.FollowerId), int(req.FollowingId))

	if err != nil {
		followHandler.logger.Error("Follow service", zap.Error(err), zap.String("requestID", requestID))
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
	return &pb.FollowResponse{
		Message: resp.Message,
	}, nil
}

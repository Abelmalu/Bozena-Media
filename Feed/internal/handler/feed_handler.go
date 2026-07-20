package handler

import (
	"context"
	"errors"

	"github.com/abelmalu/golang-posts/Feed/internal/core"
	ierrors "github.com/abelmalu/golang-posts/Feed/internal/errors"
	"github.com/abelmalu/golang-posts/Feed/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FeedHandler struct {
	pb.UnimplementedFeedServiceServer
	feedService core.FeedService
	logger      *platform.Logger
}

func NewFeedHandler(service core.FeedService, logger *platform.Logger) *FeedHandler {

	return &FeedHandler{

		feedService: service,
		logger:      logger,
	}

}

func (fh *FeedHandler) GetUserFeed(ctx context.Context, req *pb.GetUserFeedRequest) (*pb.GetUserFeedResponse, error) {

	resp, err := fh.feedService.GetUserFeed(ctx, req.Cursor, int(req.UserId), int(req.Limit))

	if err != nil {

		fh.logger.Error("Fedd Service error", zap.Error(err))

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

			return nil, status.Error(codes.DeadlineExceeded, "Request timed out")
		}

	}

	imageURL := ""

	

	pbUserFeeds := make([]*pb.UserFeed, 0, len(resp.UserFeeds))
	for _, userFeed := range resp.UserFeeds {

		if userFeed.Image == nil {

			userFeed.Image = &imageURL
		}
		pbUserFeed := &pb.UserFeed{
			PostID:      userFeed.PostID,
			PostTitle:   userFeed.PostTitle,
			PostContent: userFeed.PostContent,
			UserName:    userFeed.UserName,
			Name:        userFeed.Name,
			PostOwnerID: userFeed.PostOwnerID,
			Image: *userFeed.Image,
			LikeCount: int64(userFeed.LikeCount),
			
		}

		pbUserFeeds = append(pbUserFeeds, pbUserFeed)

	}

	return &pb.GetUserFeedResponse{
		Userfeeds: pbUserFeeds,
		Cursor:    resp.Cursor,
	}, nil

}

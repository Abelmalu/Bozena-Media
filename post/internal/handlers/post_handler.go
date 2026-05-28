package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/abelmalu/golang-posts/platform"
	"github.com/abelmalu/golang-posts/post/internal/core"
	"github.com/abelmalu/golang-posts/post/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/post/internal/errors"
	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/abelmalu/golang-posts/post/pkg/utils"
	"github.com/abelmalu/golang-posts/post/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// the PostHandler will implement the PostServiceServer
type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service core.PostService
	logger  *platform.Logger
}

func NewPostHandler(service core.PostService, logger *platform.Logger) *PostHandler {
	return &PostHandler{
		service: service,
		logger:  logger,
	}
}

func (postHandler *PostHandler) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  int(req.UserId),
	}
	validationErrStr := post.Validate()

	if validationErrStr != "" {

		return nil, status.Error(codes.InvalidArgument, validationErrStr)
	}
	var appErr *ierrors.AppError
	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}

	createdPost, err := postHandler.service.CreatePost(ctx, &post)

	if err != nil {
		postHandler.logger.Error("Post Service Error", zap.Error(err), zap.String("requestID", requestID))

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

		// Map errors to gRPC status codes
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, "request canceled")
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "timeout")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.CreatePostResponse{

		Title:   createdPost.Title,
		Content: createdPost.Content,
	}, nil

}

func (postHandler *PostHandler) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	var appErr *ierrors.AppError
	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}
	posts, err := postHandler.service.ListPosts(ctx)

	if err != nil {

		postHandler.logger.Error("Post Service Error", zap.Error(err), zap.String("requestID", requestID))

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

			return nil, status.Error(codes.Canceled, "Request Canceled")
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, status.Error(codes.DeadlineExceeded, "Request timeout")

		}
	}
	pbPosts := make([]*pb.Post, len(posts))

	for i, p := range posts {
		pbPosts[i] = &pb.Post{

			Title:   p.Title,
			Content: p.Content,
			Id:      int64(p.ID),
		}
	}

	return &pb.ListPostsResponse{
		Posts: pbPosts,
	}, nil
}

func (postHandler *PostHandler) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	post := dto.UpdatePostRequest{
		Title:   req.Title,
		Content: req.Content,
	
	}
	validationErrStr := post.Validate()

	if validationErrStr != "" {

		return nil, status.Error(codes.InvalidArgument, validationErrStr)
	}
	postID := int(req.PostId)
	var appErr *ierrors.AppError
	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}

	_, err = postHandler.service.UpdatePost(ctx, postID, post.Title, post.Content)

	if err != nil {
		postHandler.logger.Error("Post service error", zap.Error(err), zap.String("requestID", requestID))
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

			return nil, status.Error(codes.Canceled, "Request Canceled")
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, status.Error(codes.DeadlineExceeded, "Requset timeout")

		}

	}

	return &pb.UpdatePostResponse{
		Status:  "Success",
		Message: "Post updated successfully",
	}, nil

}

func (postHandler *PostHandler) GetUserPosts(ctx context.Context, req *pb.GetUserPostRequest) (*pb.GetUserPostResponse, error) {
	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {
		postHandler.logger.Error("meta data not found",zap.Error(err))
		return nil, status.Error(codes.Internal, "something went wrong")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {
		postHandler.logger.Error("requestID not found",zap.Error(err))

		return nil, status.Error(codes.InvalidArgument, "something went wrong")

	}
    resp,err := postHandler.service.GetUserPosts(ctx,req.UserId,req.Limit)

	var appErr *ierrors.AppError
	if err != nil {
		postHandler.logger.Error("Post service error", zap.Error(err), zap.String("requestID", requestID))

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

			return nil, status.Error(codes.Canceled, "Request Canceled")
		}
		if errors.Is(err, context.DeadlineExceeded) {

			return nil, status.Error(codes.DeadlineExceeded, "Requset timeout")

		}

	}
	
	pbPosts := make([]*pb.Post,len(*resp.Posts))


	for _,p := range *resp.Posts {
		pbPost := &pb.Post{

			Title:   p.Title,
			Content: p.Content,
			Id:      int64(p.ID),
			UserId:  int64(p.UserID),

		}
		pbPosts = append(pbPosts,pbPost)

	}
	
	
	return &pb.GetUserPostResponse{
		Posts:pbPosts,
		Cursor:  strconv.FormatInt(int64(resp.Cursor),16),
		HasNext: resp.HasNext,
	}, nil
}

func (postHandler *PostHandler) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	var appErr *ierrors.AppError

	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "something went wrong")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "something went wrong")

	}

	err = postHandler.service.DeletePost(ctx, int(req.PostId))
	if err != nil {

		postHandler.logger.Error("Post service error", zap.Error(err), zap.String("requestID", requestID))

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

	return &pb.DeletePostResponse{
		Status:  "Success",
		Message: "Successfully Deleted a Post",
	}, nil

}

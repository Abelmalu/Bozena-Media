package handler

import (
	"context"
	"errors"

	"github.com/abelmalu/golang-posts/Auth/internal/core"
	ierrors "github.com/abelmalu/golang-posts/Auth/internal/errors"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
	"github.com/abelmalu/golang-posts/Auth/pkg/utils"
	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	service core.AuthService
	logger  *platform.Logger
}

var validate = validator.New()

func NewAuthHandler(authServ core.AuthService, logger *platform.Logger) *AuthHandler {

	return &AuthHandler{
		service: authServ,
		logger:  logger,
	}

}

// Register registers a new user into the system
func (authHandler *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	user := model.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	validationErrStr := user.Validate()

	if validationErrStr != "" {

		return nil, status.Error(codes.InvalidArgument, validationErrStr)

	}
	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}

	createdUser, tokens, err := authHandler.service.Register(ctx, &user)
	if err != nil {
		authHandler.logger.Error("[Auth Service]", zap.Error(err), zap.String("requestID", requestID))

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

	return &pb.RegisterResponse{
		Name:         createdUser.Name,
		Username:     createdUser.Username,
		Email:        createdUser.Email,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil

}

func (authHandler *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}

	user, tokens, err := authHandler.service.Login(ctx, req.Username, req.Password)

	if err != nil {

		authHandler.logger.Error("[Auth Service]", zap.Error(err), zap.String("requestID", requestID))

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

	return &pb.LoginResponse{
		Id:             int64(user.ID),
		Username:       user.Username,
		AccessToken:    tokens.AccessToken,
		RefreshToken:   tokens.RefreshToken,
		FollowerCount:  int64(user.FollowerCount),
		FollowingCount: int64(user.FollowingCount),
	}, nil
}
func (authHandler *AuthHandler) Logout(ctx context.Context, req *emptypb.Empty) (*pb.LogoutResponse, error) {
	var refreshToken string
	var requestID string
	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")
	}
	values := md.Get("refreshToken")
	if len(values) > 0 {
		refreshToken = values[0]
	} else {

		return nil, status.Error(codes.Unauthenticated, "You are not authenticated")

	}
	requestIDValues := md.Get("requestID")
	if len(requestIDValues) > 0 {

		requestID = requestIDValues[0]
	} else {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}

	err := authHandler.service.Logout(ctx, refreshToken)

	if err != nil {

		authHandler.logger.Error("[Auth Service]", zap.Error(err), zap.String("requestID", requestID))

		var appErr *ierrors.AppError
		if errors.As(err, &appErr) {
			switch appErr.Type {
			case ierrors.TypeValidation:
				return nil, status.Error(codes.InvalidArgument, string(appErr.Message))
			case ierrors.TypeConflict:
				return nil, status.Error(codes.AlreadyExists, string(appErr.Message))
			case ierrors.TypeUnauthorized:
				return nil, status.Error(codes.Unauthenticated, string(ierrors.MSGUnauthorizedAccess))
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

	return &pb.LogoutResponse{
		Message: "Successfully Logged Out!",
	}, nil

}

func (authHandler *AuthHandler) RefreshHandler(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {

	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}

	tokens, err := authHandler.service.RefreshHandler(ctx, req.RefreshToken)

	if err != nil {

		authHandler.logger.Error("[Auth Service]", zap.Error(err), zap.String("requestID", requestID))

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

	return &pb.RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil

}

func (authHandler *AuthHandler) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {

	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}

	resp, err := authHandler.service.SearchUser(ctx, req.Username, req.Cursor, int(req.Limit))

	if err != nil {

		authHandler.logger.Error("[Auth Service]", zap.Error(err), zap.String("requestID", requestID))

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

	pbUsers := make([]*pb.User, 0, len(resp.Users))

	for _, user := range resp.Users {

		var pbUser pb.User

		pbUser = pb.User{

			Id:       int64(user.ID),
			Name:     user.Name,
			Username: user.Username,
		}

		pbUsers = append(pbUsers, &pbUser)
	}

	return &pb.SearchUserResponse{
		Users:   pbUsers,
		Cursor:  resp.Cursor,
		HasNext: resp.HasNext,
	}, nil
}

func (authHandler *AuthHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {

	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}

	resp, err := authHandler.service.GetUserProfile(ctx, req.UserId)

	if err != nil {

		authHandler.logger.Error("[Auth Service]", zap.Error(err), zap.String("requestID", requestID))

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

	avatarUrl := ""
	if resp.Avatar != nil {
		avatarUrl = *resp.Avatar
	}

	return &pb.GetUserProfileResponse{

		Id:              int64(resp.ID),
		Name:            resp.Name,
		Username:        resp.Username,
		ProfileImageUrl: avatarUrl,
		FollowerCount:   int64(resp.FollowerCount),
		FollowingCount:  int64(resp.FollowingCount),
	}, nil
}

func (authHandler *AuthHandler) GenerateUploadURL(ctx context.Context, req *pb.GenerateUploadURLRequest) (*pb.GenerateUploadURLResponse, error) {
	requestID, err := utils.GetRequestID(ctx)

	if errors.Is(err, ierrors.ErrMetaDataNotFound) {

		return nil, status.Error(codes.Internal, "meta data from context couldn't be found")

	}
	if errors.Is(err, ierrors.ErrRequestIDNotFound) {

		return nil, status.Error(codes.InvalidArgument, "missing request ID")

	}
	url, formData, err := authHandler.service.GenerateUploadURL(ctx, req.Filename, req.ContentType, int(req.UserId))

	if err != nil {

		authHandler.logger.Error("[Auth Service]", zap.Error(err), zap.String("requestID", requestID))

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
			case ierrors.TypeBadRequest:
				return nil, status.Error(codes.InvalidArgument, string(appErr.Message))
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

	return &pb.GenerateUploadURLResponse{

		UploadUrl: url,
		FormData:  formData,
	}, nil
}

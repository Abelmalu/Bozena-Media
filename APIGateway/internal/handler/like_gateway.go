package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/like/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type LikeService interface {
	ToggleLike(ctx context.Context, state bool, opts ...grpc.CallOption) (*pb.LikeResponse, error)
	GetPostLikes(ctx context.Context, postID int, limit int, cursor string) (*pb.GetPostLikesResponse, error)
}
type LikeHandler struct {
	logger     *platform.Logger
	likeClient LikeService
}

func NewLikeHandler(likeClient LikeService, logger *platform.Logger) *LikeHandler {

	return &LikeHandler{
		likeClient: likeClient,
		logger:     logger,
	}
}

// addUserIDToOutgoingContext this adds userID  and postID to the outgoing context
func addToOutgoingContext(c *gin.Context, postID string, requestID string) (context.Context, error) {

	userIDValue, exists := c.Get("userID")

	if !exists {

		return nil, ierrors.ErrUserIDNotFoundInContext
	}
	userID, ok := userIDValue.(int)

	if !ok {

		return nil, ierrors.ErrTypeAssertionFailed
	}
	userIDStr := strconv.Itoa(userID)
	md := metadata.Pairs(
		"user-id", userIDStr,
		"post-id", postID,
		"request-id", requestID,
	)

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx, nil
}
func (lh *LikeHandler) ToggleLike(c *gin.Context) {

	var req struct {
		State bool `json:"like"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lh.logger.Error("Error while Unmarshaling request body", zap.Error(err))
		c.Error(ierrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			lh.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			lh.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}

	}

	postID := c.Param("id")

	ctx, err := addToOutgoingContext(c, postID, requestID)
	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {
			lh.logger.Error("couldn't find user ID in context", zap.Error(err))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			lh.logger.Error("couldn't assert the request ID to string", zap.Error(err))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
	}
	resp, err := lh.likeClient.ToggleLike(ctx, req.State)

	if err != nil {

		lh.logger.Error("GRPC Error", zap.Error(err))
		c.Error(ierrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusCreated)
}

func (lh *LikeHandler) GetPostLikes(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			lh.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			lh.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	md := metadata.Pairs(
		"request-id", requestID,
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	followerIDStr := c.Param("id")

	postIDInt, err := strconv.Atoi(followerIDStr)

	if err != nil {

		lh.logger.Error("couldn't change followingID to string", zap.Error(err), zap.String("requestID", requestID))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return

	}

	limitStr := c.Query("limit")
	if limitStr == "" {

		limitStr = "0"
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil {

		lh.logger.Error("couldn't change followingID to string", zap.Error(err), zap.String("requestID", requestID))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return

	}
	cursor := c.Query("cursor")

	resp, err := lh.likeClient.GetPostLikes(ctx, postIDInt, limit, cursor)

	if err != nil {

		lh.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))

		c.Error(ierrors.FromGRPC(err))

		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)

}

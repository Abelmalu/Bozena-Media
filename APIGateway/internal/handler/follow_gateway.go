package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/follow/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type FollowService interface {
	ToggleFollow(ctx context.Context, follow bool, followerID, followingID int) (*pb.FollowResponse, error)
	GetUserFollowers(ctx context.Context,followingID int,limit int,cursor string)(*pb.GetUserFollowersResponse,error)
}
type FollowHandler struct {
	logger       *platform.Logger
	followClient FollowService
}

func NewFollowHandler(followClient FollowService, logger *platform.Logger) *FollowHandler {

	return &FollowHandler{
		followClient: followClient,
		logger:       logger,
	}
}

func (followHandler *FollowHandler) ToggleFollow(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			followHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			followHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	md := metadata.Pairs(
		"request-id", requestID,
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	//get userID from the context
	followerID, err := utils.GetUserID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

			followHandler.logger.Error("couldn't couldn't find userID in the context", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			followHandler.logger.Error("couldn't assert the user ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}
	followingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		followHandler.logger.Error("couldn't change followingID to string", zap.Error(err))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return

	}

	var req struct {
		Follow bool `json:"follow"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		followHandler.logger.Error("Error while unmarshalling request", zap.Error(err), zap.String("requestID", requestID))
		c.Error(ierrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return

	}
	resp, err := followHandler.followClient.ToggleFollow(ctx, req.Follow, followerID, followingID)
	if err != nil {

		followHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(ierrors.FromGRPC(err))
		return

	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusAccepted)

}


func (followHandler *FollowHandler)  GetUserFollowers(c *gin.Context) {

		requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			followHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			followHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	md := metadata.Pairs(
		"request-id", requestID,
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	followingIDStr := c.Param("id")

	followingIDInt,err := strconv.Atoi(followingIDStr)
	
	if err != nil {

		followHandler.logger.Error("couldn't change followingID to string", zap.Error(err),zap.String("requestID",requestID))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return

	}

	limitStr := c.Query("limit")

	limit,err := strconv.Atoi(limitStr)


	if err != nil {

		followHandler.logger.Error("couldn't change followingID to string", zap.Error(err),zap.String("requestID",requestID))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return

	}
	cursor := c.Query("cursor")


	resp,err := followHandler.followClient.GetUserFollowers(ctx,followingIDInt,limit,cursor)

	if err != nil {

		followHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(ierrors.FromGRPC(err))
		return

	}


	utils.SendSuccessResponse(c,resp,requestID,http.StatusOK)


}
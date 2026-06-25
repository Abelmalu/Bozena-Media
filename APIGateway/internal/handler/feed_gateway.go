package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/Feed/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)


type FeedClient interface {

	GetUserFeed(ctx context.Context,userId,limit int, cursor string)(*pb.GetUserFeedResponse,error)
}

type FeedHandler struct {

	FeedClient FeedClient 
	logger *platform.Logger

}



func NewFeedHandler(client FeedClient, logger *platform.Logger) *FeedHandler {



	return &FeedHandler{
		FeedClient:client,
		logger: logger,
	}

}


func (feedHandler *FeedHandler)  GetUserFeed(c *gin.Context) {


	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			feedHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			feedHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

		//get userID from the context
	userID, err := utils.GetUserID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

			feedHandler.logger.Error("couldn't couldn't find userID in the context", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			feedHandler.logger.Error("couldn't assert the user ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}


		limitStr := c.Query("limit")
	if limitStr == "" {

		limitStr = "0"
	}

	limit,err := strconv.Atoi(limitStr)


	if err != nil {

		feedHandler.logger.Error("couldn't change followingID to string", zap.Error(err),zap.String("requestID",requestID))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return

	}
	cursor := c.Query("cursor")

	

	md := metadata.Pairs(
		"request-id", requestID,
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp,err := feedHandler.FeedClient.GetUserFeed(ctx,userID,limit,cursor)

	if err != nil {

		feedHandler.logger.Error("GRPC Error",zap.Error(err),zap.String("requestID",requestID))

		c.Error(ierrors.FromGRPC(err))
	}


	utils.SendSuccessResponse(c,resp,requestID,http.StatusOK)



}


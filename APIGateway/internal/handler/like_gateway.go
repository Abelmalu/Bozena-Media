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
func addToOutgoingContext(c *gin.Context,postID string,requestID string)(context.Context,error){



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
		"request-id",requestID,
	
	)

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx, nil
}
func (likeHandler *LikeHandler) ToggleLike(c *gin.Context) {

	var req struct {
		State bool `json:"state"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		likeHandler.logger.Error("Error while Unmarshaling request body", zap.Error(err))
		c.Error(ierrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			likeHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			likeHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}

	}

	postID := c.Param("id")
   
	ctx, err := addToOutgoingContext(c, postID,requestID)
	if err != nil {

		if errors.Is(err,ierrors.ErrUserIDNotFoundInContext){
				likeHandler.logger.Error("couldn't find user ID in context", zap.Error(err))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return 

		}
		if errors.Is(err,ierrors.ErrTypeAssertionFailed){

				likeHandler.logger.Error("couldn't assert the request ID to string", zap.Error(err))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return 
			
		}
	}
	resp,err := likeHandler.likeClient.ToggleLike(ctx, req.State)

	if err != nil{

		likeHandler.logger.Error("GRPC Error",zap.Error(err))
		c.Error(ierrors.FromGRPC(err))
		return
	}

   utils.SendSuccessResponse(c,resp,requestID,http.StatusCreated)
}

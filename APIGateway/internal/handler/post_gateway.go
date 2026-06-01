package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	appErrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/abelmalu/golang-posts/post/proto/pb"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type PostService interface {
	CreatePost(ctx context.Context, userID int64, title, content string) (*pb.CreatePostResponse, error)
	ListPosts(ctx context.Context) (*pb.ListPostsResponse, error)
	UpdatePost(ctx context.Context, postID int64, title string, content string) (*pb.UpdatePostResponse, error)
	DeletePost(ctx context.Context, postID int64) (*pb.DeletePostResponse, error)
	GetUserPosts(ctx context.Context, userID,limit int64,cursor string) (*pb.GetUserPostResponse, error)
}

type PostHandler struct {
	postClient PostService
	logger     *platform.Logger
}

func NewPostHandler(pc PostService, logger *platform.Logger) *PostHandler {
	return &PostHandler{
		postClient: pc,
		logger:     logger,
	}
}

// this appends data to the context so grpc services can get it
func appendToOutgoingContext(c *gin.Context, requestID string) (context.Context, error) {
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
		"request-id", requestID,
	)

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx, nil
}

func (postHandler *PostHandler) CreatePost(c *gin.Context) {

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			postHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			postHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}
	//get userID from the context
	userIDInt, err := utils.GetUserID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

			postHandler.logger.Error("couldn't couldn't find userID in the context", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			postHandler.logger.Error("couldn't assert the user ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}

	if err := c.ShouldBindJSON(&req); err != nil {
		postHandler.logger.Error("Error while unmarshaling json", zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(appErrors.MSGInvalidRequestBody, nil, err))
		return
	}

	userID := int64(userIDInt)
	ctx, err := appendToOutgoingContext(c, requestID)
	if err != nil {
		postHandler.logger.Error("failed to get userID from context", zap.Error(err))
		c.Error(appErrors.NewInternalError(ierrors.MSGUnauthorizedAccess, err))
		return
	}

	resp, err := postHandler.postClient.CreatePost(ctx, userID, req.Title, req.Content)
	if err != nil {
		postHandler.logger.Error("GRPC Error", zap.Error(err))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)
}

func (postHandler *PostHandler) UpdatePost(c *gin.Context) {

	var input struct {
		Title   string `json:"title" db:"title"`
		Content string `json:"content" db:"content"`
	}
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			postHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			postHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}

	postIDStr := c.Param("id")
	postIDValue, err := strconv.Atoi(postIDStr)
	if err != nil {

		postHandler.logger.Error("error while Atoi", zap.String("requestID", requestID))

		c.Error(appErrors.NewAppError(ierrors.TypeValidation, ierrors.MSGInvalidRequestBody, err))
		return

	}
	postID := int64(postIDValue)

	if err := c.ShouldBindJSON(&input); err != nil {

		postHandler.logger.Error(string(ierrors.MSGInvalidRequestBody), zap.Error(err), zap.String("requestID", requestID))

		c.Error(ierrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}

	ctx, err := appendToOutgoingContext(c, requestID)
	if err != nil {

		postHandler.logger.Error("failed to get userID from context", zap.Error(err))
		c.Error(appErrors.NewInternalError(ierrors.MSGUnauthorizedAccess, err))
		return
	}

	resp, err := postHandler.postClient.UpdatePost(ctx, postID, input.Title, input.Content)
	if err != nil {

		postHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)

}
func (postHandler *PostHandler) DeletePost(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			postHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			postHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)

	if err != nil {

		postHandler.logger.Error("error while Atoi", zap.String("requestID", requestID))

		c.Error(appErrors.NewAppError(ierrors.TypeValidation, ierrors.MSGInvalidRequestBody, err))

		return

	}

	ctx, err := appendToOutgoingContext(c, requestID)
	if err != nil {

		postHandler.logger.Error("failed to get userID from context", zap.Error(err))
		c.Error(appErrors.NewInternalError(ierrors.MSGUnauthorizedAccess, err))
		return
	}

	resp, err := postHandler.postClient.DeletePost(ctx, int64(postID))

	if err != nil {

		postHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)

}

func (postHandler *PostHandler) GetUserPosts(c *gin.Context) {
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			postHandler.logger.Error("couldn't get request ID from context", zap.Error(err))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			postHandler.logger.Error("couldn't assert the request ID to string", zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}

	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {

		postHandler.logger.Error("Error while strConv id param", zap.Error(err))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return
	}
    
	limitStr := c.Query("limit")
	
	if limitStr == "" {

		limitStr= "0"
	}
	limit, err := strconv.Atoi(limitStr)



	if err != nil {

		postHandler.logger.Error("Error while strConv limit query param", zap.Error(err),zap.String("requestID",requestID))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return
	}

	cursor := c.Query("cursor")

	

	ctx,_ := addToOutgoingContext(c,"",requestID)

	resp, err := postHandler.postClient.GetUserPosts(ctx, int64(userID),int64(limit),cursor)
	if err != nil {

		postHandler.logger.Error("GRPC Error", zap.Error(err))

		c.Error(ierrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)

}

func (postHandler *PostHandler) ListPosts(c *gin.Context) {
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			postHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			postHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

		}
		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

		}
	}

	ctx, err := appendToOutgoingContext(c, requestID)
	if err != nil {

		postHandler.logger.Error("Failed to prepare request context", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewInternalError(ierrors.MSGUnauthorizedAccess, err))
		return
	}

	resp, err := postHandler.postClient.ListPosts(ctx)

	if err != nil {

		postHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)

}

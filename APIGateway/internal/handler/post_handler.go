package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	appErrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/pkg/utils"
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

func addUserIDToContext(c *gin.Context) (context.Context, error) {
	userIDValue, exists := c.Get("userID")

	if !exists {

		return nil, errors.New("user not found in the request")
	}
	userID, ok := userIDValue.(int)

	if !ok {

		return nil, errors.New("assertion failed on userID")
	}
	userIDStr := strconv.Itoa(userID)
	md := metadata.Pairs("user-id", userIDStr)

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx, nil
}

func (postHandler *PostHandler) CreatePost(c *gin.Context) {

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	requestID, err := utils.GetRequestID(c, postHandler.logger)

	if err != nil {

		return
	}

	//get userID from the context
	userIDInt, err := utils.GetUserID(c, postHandler.logger)
	if err != nil {

		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		postHandler.logger.Error(string(ierrors.MSGInvalidRequestBody), zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(appErrors.MSGInvalidRequestBody, nil, err))
		return
	}

	userID := int64(userIDInt)
	ctx, err := addUserIDToContext(c)
	if err != nil {
		postHandler.logger.Error("", zap.Error(err))
		c.Error(appErrors.NewInternalError("Failed to prepare request context", err))
		return
	}

	resp, err := postHandler.postClient.CreatePost(ctx, userID, req.Title, req.Content)
	if err != nil {
		postHandler.logger.Error("GRPC Error", zap.Error(err))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (postHandler *PostHandler) UpdatePost(c *gin.Context) {

	var input struct {
		Title   string `json:"title" db:"title"`
		Content string `json:"content" db:"content"`
	}
	requestID, err := utils.GetRequestID(c, postHandler.logger)
	if err != nil {

		return
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

	ctx, err := addUserIDToContext(c)
	if err != nil {

		postHandler.logger.Error(string(ierrors.MSGInvalidRequestBody), zap.Error(err), zap.String("requestID", requestID))

		c.Error(appErrors.NewInternalError("Failed to prepare request context", err))
		return
	}

	resp, err := postHandler.postClient.UpdatePost(ctx, postID, input.Title, input.Content)
	if err != nil {

		postHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	c.JSON(http.StatusCreated, resp)
}
func (postHandler *PostHandler) DeletePost(c *gin.Context) {

	requestID, err := utils.GetRequestID(c, postHandler.logger)
	if err != nil {

		return
	}
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)

	if err != nil {

		postHandler.logger.Error("error while Atoi", zap.String("requestID", requestID))

		c.Error(appErrors.NewAppError(ierrors.TypeValidation, ierrors.MSGInvalidRequestBody, err))

		return

	}

	ctx, err := addUserIDToContext(c)
	if err != nil {

		postHandler.logger.Error(string(ierrors.MSGInvalidRequestBody), zap.Error(err), zap.String("requestID", requestID))

		c.Error(appErrors.NewInternalError("Failed to prepare request context", err))
		return
	}

	resp, err := postHandler.postClient.DeletePost(ctx, int64(postID))

	if err != nil {

		postHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	c.JSON(http.StatusAccepted, resp)
}

func (postHandler *PostHandler) ListPosts(c *gin.Context) {
	requestID,err := utils.GetRequestID(c,postHandler.logger)
	if err != nil{

		return
	}

	ctx, err := addUserIDToContext(c)
	if err != nil {

		postHandler.logger.Error(string(ierrors.MSGInvalidRequestBody), zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewInternalError("Failed to prepare request context", err))
		return
	}

	resp, err := postHandler.postClient.ListPosts(ctx)

	if err != nil {

		postHandler.logger.Error("GRPC Error",zap.Error(err),zap.String("requestID",requestID))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	c.JSON(http.StatusOK, resp)

}

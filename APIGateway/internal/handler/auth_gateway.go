package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/abelmalu/golang-posts/APIGateway/internal/dto"
	appErrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type AuthService interface {
	Register(ctx context.Context, userName, name, email, password string) (*pb.RegisterResponse, error)
	Login(ctx context.Context, userName, password string) (*pb.LoginResponse, error)
	Logout(ctx context.Context) (*pb.LogoutResponse, error)
	RefreshHandler(context.Context, string) (*pb.RefreshResponse, error)
	GetUserProfile(ctx context.Context, userID int64) (*pb.GetUserProfileResponse, error)
	SearchUser(ctx context.Context, username, cursor string, limit int) (*pb.SearchUserResponse, error)
	GenerateProfileUploadURL(ctx context.Context, userID int, fileName, ContentType string) (*pb.GenerateUploadURLResponse, error)
}
type AuthHandler struct {
	client AuthService
	logger *platform.Logger
}

func NewAuthHandler(au AuthService, logger *platform.Logger) *AuthHandler {

	return &AuthHandler{
		client: au,
		logger: logger,
	}
}

// ExtractRefreshToken extracts refresh tokens from the request
func ExtractRefreshToken(c *gin.Context) (string, error) {

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	//for mobile apps from request body
	if err := c.ShouldBindJSON(&req); err == nil {
		if req.RefreshToken != "" {
			return req.RefreshToken, nil
		}
	}

	// for browsers
	if token, err := c.Cookie("refresh_token"); err == nil && token != "" {
		return token, nil
	}

	return "", errors.New("Refresh Token not found")
}

func (authHandler *AuthHandler) Register(c *gin.Context) {

	var req struct {
		Name     string `json:"name"`
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			authHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	if err := c.ShouldBindJSON(&req); err != nil {
		authHandler.logger.Error("error while marshaling request", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}
	// call getClienType to get the client type and inject it into the grpc metadata
	ctx, clientType := utils.AddToOutgoingContext(c, requestID)

	resp, err := authHandler.client.Register(ctx, req.UserName, req.Name, req.Email, req.Password)
	if err != nil {

		authHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	var registerResponse dto.RegisterResponse

	switch clientType {
	case "web":
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			MaxAge:   30 * 24 * 60 * 60,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		registerResponse.AccessToken = resp.AccessToken

	case "mobile":
		registerResponse.AccessToken = resp.AccessToken
		registerResponse.RefreshToken = resp.RefreshToken
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)

}

func (authHandler *AuthHandler) Login(c *gin.Context) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			authHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}
	if err := c.ShouldBindJSON(&req); err != nil {
		authHandler.logger.Error("error while marshaling request", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}

	ctx, clientType := utils.AddToOutgoingContext(c, requestID)

	resp, err := authHandler.client.Login(ctx, req.UserName, req.Password)

	if err != nil {
		authHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("request ID", requestID))
		c.Error(appErrors.FromGRPC(err))
		return
	}
	var loginResponse dto.LoginResponse

	switch clientType {
	case "web":
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			MaxAge:   30 * 24 * 60 * 60,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		loginResponse.AccessToken = resp.AccessToken
		loginResponse.UserName = resp.Username
		loginResponse.FollowerCount = int(resp.FollowerCount)
		loginResponse.FollowingCount = int(resp.FollowingCount)
		loginResponse.ID = int(resp.Id)
		loginResponse.ProfileImageUrl =resp.ProfileImageUrl 

	case "mobile":
		loginResponse.AccessToken = resp.AccessToken
		loginResponse.RefreshToken = resp.RefreshToken
		loginResponse.UserName = resp.Username
		loginResponse.FollowerCount = int(resp.FollowerCount)
		loginResponse.FollowingCount = int(resp.FollowingCount)
		loginResponse.ID = int(resp.Id)
		loginResponse.ProfileImageUrl =resp.ProfileImageUrl 

	}

	utils.SendSuccessResponse(c, loginResponse, requestID, http.StatusOK)

}

func (authHandler *AuthHandler) Logout(c *gin.Context) {
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			authHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}
	}

	JTI, err := utils.GetJTI(c)

	if err != nil {

		if errors.Is(err, ierrors.ErrJTINotFoundInContext) {

			authHandler.logger.Error("couldn't get JTI from context", zap.Error(errors.New("couldn't find JTI")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}

		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the JTI to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	expTime, err := utils.GetJWTEXPTime(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrExpTimeNotFoundInContext) {

			authHandler.logger.Error("couldn't get expirtaion time from context", zap.Error(errors.New("couldn't find expiration time")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}

		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}
	expTimeStr := strconv.Itoa(int(expTime))

	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {
		authHandler.logger.Error("refresh token extracting error", zap.Error(err), zap.String("request ID", requestID))
		c.Error(appErrors.NewAppError(appErrors.TypeUnauthorized, ierrors.MSGRefreshTokenNotFound, err))
		return
	}

	md := metadata.Pairs(
		"refreshToken", refreshToken,
		"requestID", requestID,
		"JTI", JTI,
		"expTime", expTimeStr,
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)
	resp, err := authHandler.client.Logout(ctx)
	if err != nil {

		authHandler.logger.Error("GRPC Error", zap.Error(err))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)
}

func (authHandler *AuthHandler) RefreshHandler(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			authHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	// extracting the refresh token from the request for both mobile and web clients
	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {
		authHandler.logger.Error("refresh token extracting error ", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewAppError(appErrors.TypeUnauthorized, "Refresh token not found", err))
		return
	}

	ctx, clientType := utils.AddToOutgoingContext(c, requestID)

	resp, err := authHandler.client.RefreshHandler(ctx, refreshToken)
	if err != nil {
		authHandler.logger.Error("GRPC Error ", zap.Error(err), zap.String("requestID", requestID))

		c.Error(appErrors.FromGRPC(err))
		return
	}

	response := gin.H{"message": "Refreshed successfully"}

	switch clientType {
	case "web":
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			MaxAge:   30 * 24 * 60 * 60,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		response["access_token"] = resp.AccessToken

	case "mobile":
		response["access_token"] = resp.AccessToken
		response["refresh_token"] = resp.RefreshToken
	}

	c.JSON(http.StatusOK, response)

}

func (authHandler *AuthHandler) SearchUser(c *gin.Context) {
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			authHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	search := c.Query("search")

	limitStr := c.Query("limit")

	if limitStr == "" {

		limitStr = "0"
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil {

		authHandler.logger.Error("couldn't change limit to string", zap.Error(err), zap.String("requestID", requestID))
		c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
		return

	}

	cursor := c.Query("cursor")

	ctx, _ := utils.AddToOutgoingContext(c, requestID)

	resp, err := authHandler.client.SearchUser(ctx, search, cursor, limit)

	if err != nil {

		authHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestId", requestID))

		c.Error(ierrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)
}

func (authHandler *AuthHandler) GetUserProfile(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			authHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	userIDStr := c.Param("id")
	userIDValue, err := strconv.Atoi(userIDStr)
	if err != nil {

		authHandler.logger.Error("error while Atoi", zap.String("requestID", requestID))

		c.Error(appErrors.NewAppError(ierrors.TypeValidation, ierrors.MSGInvalidRequestBody, err))
		return

	}

	ctx, _ := utils.AddToOutgoingContext(c, requestID)

	userID := int64(userIDValue)

	resp, err := authHandler.client.GetUserProfile(ctx, userID)

	if err != nil {

		authHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(ierrors.FromGRPC(err))

		return

	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)
}

func (authHandler *AuthHandler) GenerateProfileUploadURL(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			authHandler.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}
	userID, err := utils.GetUserID(c)

	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

			authHandler.logger.Error("couldn't couldn't find userID in the context", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			authHandler.logger.Error("couldn't assert the user ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}

	var req struct {
		FileName    string `json:"file_name"`
		ContentType string `json:"content_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		authHandler.logger.Error("error while marshaling request", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return

	}
	ctx, _ := utils.AddToOutgoingContext(c, requestID)

	resp, err := authHandler.client.GenerateProfileUploadURL(ctx, userID, req.FileName, req.ContentType)

	if err != nil {

		authHandler.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
		c.Error(ierrors.FromGRPC(err))

		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)

}

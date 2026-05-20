package handler

import (
	"context"
	"errors"
	"net/http"

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

// getClientType get client type header and inject into the contex metadata
func getClientType(c *gin.Context, requestID string) (context.Context, string) {

	clientType := c.GetHeader("X-Client-Type")
	md := metadata.Pairs(
		"x-client-type", clientType,
		"requestID", requestID,
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx, clientType

}

// ExtractRefreshToken extracts refresh tokens from the request
func ExtractRefreshToken(c *gin.Context) (string, error) {

	var refreshToken string

	//for mobile apps from request body
	if err := c.ShouldBindJSON(&refreshToken); err == nil {
		if refreshToken != "" {
			return refreshToken, nil
		}
	}

	// for browsers
	if token, err := c.Cookie("refresh_token"); err == nil && token != "" {
		return token, nil
	}



	return "", errors.New("Refresh Token not found")
}

func (ah *AuthHandler) Register(c *gin.Context) {

	var req struct {
		Name     string `json:"name"`
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	requestID,_ := utils.GetRequestID(c, ah.logger)

	if err := c.ShouldBindJSON(&req); err != nil {
		ah.logger.Error("error while marshaling request", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}
	// call getClienType to get the client type and inject it into the grpc metadata
	ctx, clientType := getClientType(c, requestID)

	resp, err := ah.client.Register(ctx, req.UserName, req.Name, req.Email, req.Password)
	if err != nil {

		ah.logger.Error("GRPC Error", zap.Error(err), zap.String("requestID", requestID))
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

	utils.SendSuccessResponse(c,registerResponse,requestID,http.StatusOK)

}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	requestID,_ := utils.GetRequestID(c, ah.logger)

	if err := c.ShouldBindJSON(&req); err != nil {
		ah.logger.Error("Error Unmarshaling request", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}
	ctx, clientType := getClientType(c, requestID)

	resp, err := ah.client.Login(ctx, req.UserName, req.Password)

	if err != nil {
		ah.logger.Error("GRPC Error", zap.Error(err), zap.String("request ID", requestID))
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

		loginResponse.AccessToken= resp.AccessToken

	case "mobile":
		loginResponse.AccessToken = resp.AccessToken
		loginResponse.RefreshToken = resp.RefreshToken
	}

	utils.SendSuccessResponse(c,loginResponse,requestID,http.StatusOK)

}

func (ah *AuthHandler) Logout(c *gin.Context) {
	requestID,_ := utils.GetRequestID(c, ah.logger)

	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {
		ah.logger.Error("refresh token extracting error", zap.Error(err), zap.String("request ID", requestID))
		c.Error(appErrors.NewAppError(appErrors.TypeUnauthorized, ierrors.MSGRefreshTokenNotFound, err))
		return
	}
	md := metadata.Pairs("refreshToken", refreshToken, "requestID", requestID)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)
	resp, err := ah.client.Logout(ctx)
	if err != nil {

		ah.logger.Error("GRPC Error", zap.Error(err))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c,resp,requestID,http.StatusOK)
}

func (ah *AuthHandler) RefreshHandler(c *gin.Context) {

	requestID,_ := utils.GetRequestID(c, ah.logger)

	// extracting the refresh token from the request for both mobile and web clients
	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {
		ah.logger.Error("refresh token extracting error ", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewAppError(appErrors.TypeUnauthorized, "Refresh token not found", err))
		return
	}

	ctx, clientType := getClientType(c, requestID)

	resp, err := ah.client.RefreshHandler(ctx, refreshToken)
	if err != nil {
		ah.logger.Error("GRPC Error ", zap.Error(err), zap.String("requestID", requestID))

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

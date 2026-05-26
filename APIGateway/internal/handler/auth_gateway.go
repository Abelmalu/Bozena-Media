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
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			ah.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			ah.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ah.logger.Error("error while marshaling request", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}
	// call getClienType to get the client type and inject it into the grpc metadata
	ctx, clientType := utils.AddToOutgoingContext(c, requestID)

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

	utils.SendSuccessResponse(c, registerResponse, requestID, http.StatusOK)

}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			ah.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			ah.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}
    if err := c.ShouldBindJSON(&req); err != nil {
		ah.logger.Error("error while marshaling request", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewValidationError(ierrors.MSGInvalidRequestBody, nil, err))
		return
	}
	
	ctx, clientType := utils.AddToOutgoingContext(c, requestID)

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

		loginResponse.AccessToken = resp.AccessToken

	case "mobile":
		loginResponse.AccessToken = resp.AccessToken
		loginResponse.RefreshToken = resp.RefreshToken
	}

	utils.SendSuccessResponse(c, loginResponse, requestID, http.StatusOK)

}

func (ah *AuthHandler) Logout(c *gin.Context) {
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			ah.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			ah.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}
	}
    

	JTI,err :=utils.GetJTI(c)

	if err != nil {


		if errors.Is(err, ierrors.ErrJTINotFoundInContext) {

			ah.logger.Error("couldn't get JTI from context", zap.Error(errors.New("couldn't find JTI")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}

		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			ah.logger.Error("couldn't assert the JTI to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}




	}

	expTime,err := utils.GetJWTEXPTime(c)
	if err != nil {


		if errors.Is(err, ierrors.ErrExpTimeNotFoundInContext) {

			ah.logger.Error("couldn't get expirtaion time from context", zap.Error(errors.New("couldn't find expiration time")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}

		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			ah.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}




	}
	expTimeStr := strconv.Itoa(int(expTime))

	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {
		ah.logger.Error("refresh token extracting error", zap.Error(err), zap.String("request ID", requestID))
		c.Error(appErrors.NewAppError(appErrors.TypeUnauthorized, ierrors.MSGRefreshTokenNotFound, err))
		return
	}

	md := metadata.Pairs(
		"refreshToken", refreshToken, 
		 "requestID",requestID,
		 "JTI",JTI,
		 "expTime",expTimeStr,

		
		)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)
	resp, err := ah.client.Logout(ctx)
	if err != nil {

		ah.logger.Error("GRPC Error", zap.Error(err))
		c.Error(appErrors.FromGRPC(err))
		return
	}

	utils.SendSuccessResponse(c, resp, requestID, http.StatusOK)
}

func (ah *AuthHandler) RefreshHandler(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			ah.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			ah.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}


	}

	// extracting the refresh token from the request for both mobile and web clients
	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {
		ah.logger.Error("refresh token extracting error ", zap.Error(err), zap.String("requestID", requestID))
		c.Error(appErrors.NewAppError(appErrors.TypeUnauthorized, "Refresh token not found", err))
		return
	}

	ctx, clientType := utils.AddToOutgoingContext(c, requestID)

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

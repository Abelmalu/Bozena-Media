package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type AuthService interface {
	Register(ctx context.Context, userName, name, email, password string) (*pb.RegisterResponse, error)
	Login(ctx context.Context, userName, password string) (*pb.LoginResponse, error)
	Logout(ctx context.Context) (*pb.LogoutResponse, error)
	RefreshHandler(context.Context,string)(*pb.RefreshResponse,error)
}
type AuthHandler struct {
	client AuthService
}

func NewAuthHandler(au AuthService) *AuthHandler {

	return &AuthHandler{client: au}
}

// getClientType get client type header and inject into the contex metadata
func getClientType(c *gin.Context) (context.Context, string) {

	clientType := c.GetHeader("X-Client-Type")
	md := metadata.Pairs("x-client-type", clientType)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx, clientType

}

// ExtractRefreshToken extracts refresh tokens from the request 
func ExtractRefreshToken(c *gin.Context) (string, error) {

	//Extracting  refresh tokens from mobile app clients
	var refreshToken string
	if err := c.ShouldBindJSON(&refreshToken); err == nil {
		if refreshToken != "" {
			return refreshToken, nil
		}
	}

	// Extracting HttpOnly cookie for web clients
	if token, err := c.Cookie("refresh_token"); err == nil && token != "" {
		return token, nil
	}

	//Extracting from custom header fallback
	if token := c.GetHeader("X-Refresh-Token"); token != "" {
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
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf(" error while decoding json %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
		})

		return

	}
	// call getClienType to get the client type and inject it into the grpc metadata
	ctx, clientType := getClientType(c)

	resp, err := ah.client.Register(ctx, req.UserName, req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("the error while calling client service %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}
	response := gin.H{"message": "Registered successfully"}

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

func (ah *AuthHandler) Login(c *gin.Context) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
		})

		return

	}
	ctx, clientType := getClientType(c)

	resp, err := ah.client.Login(ctx, req.UserName, req.Password)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return

	}
	response := gin.H{"message": "Registered successfully"}

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

func (ah *AuthHandler) Logout(c *gin.Context){
	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {

		log.Printf("refresh token extracting error %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return
	}
	md := metadata.Pairs("refreshToken", refreshToken)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)
	resp,err := ah.client.Logout(ctx)

	c.JSON(http.StatusAccepted,resp)
}


func(ah *AuthHandler) RefreshHandler(c *gin.Context) {

	// extracting the refresh token from the request for both mobile and web clients
	refreshToken, err := ExtractRefreshToken(c)
	if err != nil {

		log.Printf("refresh token extracting error %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return
	}

	ctx,clientType := getClientType(c)

	resp,err := ah.client.RefreshHandler(ctx,refreshToken)
	if err != nil {

		log.Printf("refresh token  error %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return


	}

	response := gin.H{"message": "Registered successfully"}

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

	c.JSON(http.StatusOK, resp)



}

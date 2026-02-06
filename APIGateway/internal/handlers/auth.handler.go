package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type AuthService interface {
	Register(ctx context.Context, userName, name, email, password string) (*pb.RegisterResponse, error)
	Login(ctx context.Context, userName, password string) (*pb.LoginResponse, error)
	Logout(ctx context.Context, userName, password string) (*pb.LogoutResponse, error)
}
type AuthHandler struct {
	client AuthService
}

func NewAuthHandler(au AuthService) *AuthHandler {

	return &AuthHandler{client: au}
}



func (ah *AuthHandler) Register(c *gin.Context) {

	var req struct {
		Name     string `json:"name"`
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf(" error while decoding json %v",err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
		})

		return

	}
	clientTypeValue := c.GetHeader("X-Client-Type")
	// if !exists {
	// 	log.Printf("here in the custom header %v",exists)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "error",
	// 		"message": "Bad Request",
	// 	})

	// 	return

	// }
	//type insertion for clientype string 
	clientType := clientTypeValue
  
	md := metadata.Pairs("x-client-type",clientType)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp, err := ah.client.Register(ctx, req.UserName, req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("the error while calling client service %v",err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}
	response := gin.H{"message": "Registered successfully"}

	switch clientType{
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

	resp, err := ah.client.Login(c, req.UserName, req.Password)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusAccepted, resp)

}

// func (ah *AuthHandler) Logout(c *gin.Context){
// 	resp,err := ah.client.Logout()
// }

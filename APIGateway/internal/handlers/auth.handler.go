package handler

import (
	"context"
	"net/http"

	"github.com/abelmalu/golang-posts/Auth/proto/pb"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
		})

		return

	}
	resp,err := ah.client.Register(c.Request.Context(),req.UserName,req.Name,req.Email,req.Password)
	if err != nil{

		c.JSON(http.StatusInternalServerError,gin.H{
			"status":"error",
			"message":"Internal Server Error",
		})
		return 
	}
	c.JSON(http.StatusCreated,resp)

}

func (ah *AuthHandler) Login(c *gin.Context){
	var req struct{
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

	resp,err := ah.client.Login(c,req.UserName,req.Password)

	if err != nil{

		c.JSON(http.StatusInternalServerError,gin.H{
			"status":"error",
			"message":"Internal Server Error",
		})
		return 
	}

	c.JSON(http.StatusAccepted,resp)

}


// func (ah *AuthHandler) Logout(c *gin.Context){
// 	resp,err := ah.client.Logout()
// }
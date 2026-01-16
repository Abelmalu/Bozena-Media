package internal

import (
	"github.com/abelmalu/golang-posts/internal/middleware"
	"github.com/abelmalu/golang-posts/internal/posts"
	"github.com/abelmalu/golang-posts/internal/auth"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	r := gin.Default()
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/refresh",auth.RefreshHandler)

	}
	postsGroup := r.Group("/posts")
	postsGroup.Use(middleware.AuthMiddleware())
	{
		postsGroup.POST("/create", posts.CreatePost)
	}

	return r
}

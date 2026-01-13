package internal

import (
	"github.com/abelmalu/golang-posts/internal/auth"
	"github.com/abelmalu/golang-posts/internal/middleware"
	"github.com/abelmalu/golang-posts/internal/posts"
	"github.com/gin-gonic/gin"

)

func SetupRoutes() *gin.Engine {

	r := gin.Default()
	postsGroup := r.Group("/posts")
	postsGroup.Use(middleware.AuthMiddleware())
	{
		postsGroup.POST("/create", posts.CreatePost)
	}

	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	

	return r
}


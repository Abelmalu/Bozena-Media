package internal

import (
	"github.com/abelmalu/golang-posts/internal/middleware"
	"github.com/abelmalu/golang-posts/internal/posts"
	"github.com/abelmalu/golang-posts/internal/auth"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/gin-contrib/cors"
)

func SetupRoutes() *gin.Engine {

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization","X-Client-Type"},
		ExposeHeaders:    []string{"Authorization"}, // IMPORTANT
		AllowCredentials: true, // REQUIRED for cookies
		MaxAge:           12 * time.Hour,
	}))

	authGroup := r.Group("")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/refresh",auth.RefreshHandler)
		authGroup.POST("/logout",middleware.AuthMiddleware(),auth.Logout)
		
	}
	postsGroup := r.Group("/posts")
	
	postsGroup.Use(middleware.AuthMiddleware(),middleware.AuthorizeRoles("users"))
	{
		
		postsGroup.POST("/create", posts.CreatePost)
		postsGroup.GET("/",posts.GetPosts)
		postsGroup.PUT("/update/:id",posts.UpdatePost)
	}

	return r
}

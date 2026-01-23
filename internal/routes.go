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



	authGroup := r.Group("/")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/refresh",auth.RefreshHandler)
		authGroup.POST("/logout",middleware.AuthMiddleware(),auth.Logout)
		
	}
	authenticated := r.Group("/")
	
	authenticated.Use(middleware.AuthMiddleware(),middleware.AuthorizeRoles("users"))
	{
		
		authenticated.POST("/posts", posts.CreatePost)
		authenticated.GET("/posts",posts.GetPosts)

		postOwner := authenticated.Group("/posts/:id")
		postOwner.Use(middleware.AuthorizePostOwner())
		postOwner.PUT("",posts.UpdatePost)
		postOwner.DELETE("",posts.DeletePost)
	}
		

	return r
}

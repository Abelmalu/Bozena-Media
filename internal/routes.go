package internal

import (
	"time"

	"github.com/abelmalu/golang-posts/internal/auth"
	"github.com/abelmalu/golang-posts/internal/middleware"
	post "github.com/abelmalu/golang-posts/internal/posts"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Client-Type"},
		ExposeHeaders:    []string{"Authorization"}, // IMPORTANT
		AllowCredentials: true,                      // REQUIRED for cookies
		MaxAge:           12 * time.Hour,
	}))

	authGroup := r.Group("/")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/refresh", auth.RefreshHandler)
		authGroup.POST("/logout", middleware.AuthMiddleware(), auth.Logout)

	}
	authenticated := r.Group("/")

	authenticated.Use(middleware.AuthMiddleware(), middleware.AuthorizeRoles("users"))
	{

		posts := authenticated.Group("/posts")
		{
			posts.POST("", post.CreatePost)
			posts.GET("", post.GetPosts)

		}

		postOwner := posts.Group("/:id")
		{

			postOwner.Use(middleware.AuthorizePostOwner())
			postOwner.PUT("", post.UpdatePost)
			postOwner.DELETE("", post.DeletePost)

		}
	}

	return r
}

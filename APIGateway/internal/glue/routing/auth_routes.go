package routing

import (
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/glue"
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)
func InitAuthRoute(router *gin.RouterGroup, handler *handler.AuthHandler,logger *platform.Logger,redisClient *redis.Client) {


	routes := []glue.Route{
		{
			Method:      http.MethodPost,
			Path:        "/register",
			Handler:     handler.Register,
			Middlewares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/login",
			Handler:     handler.Login,
			Middlewares: [] gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/refresh",
			Handler:     handler.RefreshHandler,
			Middlewares: [] gin.HandlerFunc{},
		},
		{
			Method:  http.MethodPost,
			Path:    "/logout",
			Handler: handler.Logout,
			Middlewares: [] gin.HandlerFunc{
				middleware.AuthMiddleware(logger,redisClient),
			},
		},
			{
			Method:  http.MethodGet,
			Path:    "/search",
			Handler: handler.SearchUser,
			Middlewares: [] gin.HandlerFunc{
				middleware.AuthMiddleware(logger,redisClient),
			},
		},

		{
			Method:  http.MethodGet,
			Path:    "/profile/:id",
			Handler: handler.GetUserProfile,
			Middlewares: [] gin.HandlerFunc{
				middleware.AuthMiddleware(logger,redisClient),
			},
		},

		{
			Method:  http.MethodPost,
			Path:    "/profile/upload",
			Handler: handler.GenerateProfileUploadURL,
			Middlewares: [] gin.HandlerFunc{
				middleware.AuthMiddleware(logger,redisClient),
			},
		},
	}

	glue.RegisterRoutes(router,routes)

}

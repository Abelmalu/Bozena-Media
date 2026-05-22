package routing

import (
	"net/http"

	"github.com/abelmalu/golang-posts/APIGateway/internal/glue"
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
	"github.com/abelmalu/golang-posts/APIGateway/internal/middleware"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
)







func InitLikeRoute(router *gin.RouterGroup, handler *handler.LikeHandler,logger *platform.Logger){

	routes := []glue.Route{

		{
			Method: http.MethodPost,
			Path:"/like/:id",
			Handler: handler.ToggleLike,
			Middlewares: []func(*platform.Logger) gin.HandlerFunc{

				middleware.AuthMiddleware,
			},


		},


	}
	glue.RegisterRoutes(router,routes,logger)
}
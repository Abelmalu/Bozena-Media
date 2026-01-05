package internal

import (
	"github.com/abelmalu/golang-posts/internal/auth"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	r := gin.Default()

	r.POST("/register", auth.Register)

	return r
}

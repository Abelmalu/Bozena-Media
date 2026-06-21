package middleware

import (
	_ "embed"
	"time"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

//go:embed rate_limiter.lua
var luaScriptString string

// Initialize the script with the embedded string
var luaScript = redis.NewScript(luaScriptString)

func RateLimitMiddleware(redisClient *redis.Client,logger *platform.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "bucket:" + c.ClientIP()
		now := float64(time.Now().UnixNano()) / 1e9

		result, err := luaScript.Run(
			c.Request.Context(),
			redisClient,
			[]string{key},
			2,
			1,
			now,
			1,
		).Result()

		if err != nil {

			logger.Error("rate limiter failed",zap.Error(err))
			internalErr := ierrors.NewInternalError(ierrors.MSGSomethingWentWrong,err)
			c.Error(internalErr)
			c.Abort()
			return
		}

		values := result.([]interface{})
		allowed := values[0].(int64)

		if allowed == 0 {
			internalErr := ierrors.NewTooManyRequestsError(ierrors.MSGTooManyRequests,nil,nil)
			logger.Warn("too many requests",zap.Error(internalErr))
			c.Error(internalErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

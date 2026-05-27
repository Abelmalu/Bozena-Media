package middleware

import (
	_ "embed"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

//go:embed rate_limiter.lua
var luaScriptString string

// Initialize the script with the embedded string
var luaScript = redis.NewScript(luaScriptString)

func RateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
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
			c.JSON(500, gin.H{"error": "rate limiter failed"})
			c.Abort()
			return
		}

		values := result.([]interface{})
		allowed := values[0].(int64)

		if allowed == 0 {
			c.JSON(429, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

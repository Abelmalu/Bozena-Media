package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)
// lua programming lanuage script for atomicity in redis 
var luaScript = redis.NewScript(`
local key = KEYS[1]

local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local data = redis.call("HMGET", key, "tokens", "last_refill")

local tokens = tonumber(data[1])
local last_refill = tonumber(data[2])

if tokens == nil then
    tokens = capacity
    last_refill = now
end

local delta = math.max(0, now - last_refill)

local filled_tokens = tokens + (delta * rate)

if filled_tokens > capacity then
    filled_tokens = capacity
end

local allowed = filled_tokens >= requested

local new_tokens = filled_tokens

if allowed then
    new_tokens = filled_tokens - requested
end

redis.call(
    "HMSET",
    key,
    "tokens", new_tokens,
    "last_refill", now
)

redis.call("EXPIRE", key, 3600)

return {
    allowed and 1 or 0,
    new_tokens
}
`)

func RateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {

	return func(c *gin.Context) {
		key := "bucket:" + c.ClientIP()

		now := float64(time.Now().UnixNano()) / 1e9

		result, err := luaScript.Run(
			c.Request.Context(),
			redisClient,
			[]string{key},
			100, 
			120,  
			now,
			1,    
		).Result()

		if err != nil {

			c.JSON(500, gin.H{
				"error": "rate limiter failed",
			})
			c.Abort()
			return
		}

		values := result.([]interface{})

		allowed := values[0].(int64)

		if allowed == 0 {

			c.JSON(429, gin.H{
				"error": "rate limit exceeded",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type rateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func RateLimiter(client *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	rl := &rateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}

	return func(c *gin.Context) {
		key := c.ClientIP()
		ok, err := rl.allow(c.Request.Context(), key)
		if err != nil {
			//TODO: logs
			c.Next()
			return
		}
		if !ok {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

func (rl *rateLimiter) allow(ctx context.Context, key string) (bool, error) {
	pipe := rl.client.TxPipeline()
	inc := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rl.window)

	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}
	return inc.Val() <= int64(rl.limit), nil
}

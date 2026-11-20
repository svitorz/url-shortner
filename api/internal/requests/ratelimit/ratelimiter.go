package ratelimit

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func New(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{client: client, limit: limit, window: window}
}

func (rl *RateLimiter) Allow(c *gin.Context) (bool, error) {
	key := c.ClientIP()
	ok, err := rl.allow(c.Request.Context(), key)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (rl *RateLimiter) allow(ctx context.Context, key string) (bool, error) {
	pipe := rl.client.TxPipeline()
	inc := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rl.window)

	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}
	return inc.Val() <= int64(rl.limit), nil
}

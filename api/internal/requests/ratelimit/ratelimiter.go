package ratelimit

import (
	"context"
	"svitorz/url-shortner/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type rateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func RateLimiter(c *gin.Context) (bool, error) {
	limit := 10
	window := time.Minute
	client := repository.RDB
	rl := &rateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}

	key := c.ClientIP()
	ok, err := rl.allow(c.Request.Context(), key)
	if err != nil {
		return false, err
	}

	return ok, nil
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

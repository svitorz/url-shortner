package ratelimit

import (
	"context"
	"fmt"
	"svitorz/url-shortner/internal/config"
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

var rl *rateLimiter

func init() {
	cfg, err := config.LoadRateLimitConf()

	if err != nil {
		fmt.Println("Erro ao buscar configuracoes do rate limiter")
	}

	rl = &rateLimiter{
		client: repository.RDB,
		limit:  cfg[0],
		window: time.Duration(cfg[1] * time.Now().Minute()),
	}
}
func RateLimiter(c *gin.Context) (bool, error) {
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

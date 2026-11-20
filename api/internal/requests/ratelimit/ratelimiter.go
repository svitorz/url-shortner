package ratelimit

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// keyPrefix separa as chaves do limiter das demais (ex.: cache de links).
const keyPrefix = "ratelimit:"

// RateLimiter é uma janela fixa por IP: até limit requisições por window.
type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func New(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{client: client, limit: limit, window: window}
}

func (rl *RateLimiter) Allow(c *gin.Context) (bool, error) {
	return rl.allow(c.Request.Context(), keyPrefix+c.ClientIP())
}

func (rl *RateLimiter) allow(ctx context.Context, key string) (bool, error) {
	// MULTI/EXEC: INCR + EXPIRE NX. O NX só define o TTL se a chave ainda não
	// tem um, ou seja, na primeira requisição da janela: requisições seguintes
	// não renovam o prazo (janela fixa). Também corrige uma chave que, por
	// qualquer motivo, tenha ficado sem TTL, em vez de bloquear o IP para sempre.
	pipe := rl.client.TxPipeline()
	inc := pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, rl.window)

	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}
	return inc.Val() <= int64(rl.limit), nil
}

package repository

import (
	"svitorz/url-shortner/internal/config"
	"testing"
)

func TestRedisOptionsRespectContextDeadline(t *testing.T) {
	if !redisOptions(config.RedisConfig{}, "x:1").ContextTimeoutEnabled {
		t.Fatal("ContextTimeoutEnabled desligado: o timeout do cache de links seria ignorado")
	}
}

package repository

import (
	"context"
	"fmt"
	"svitorz/url-shortner/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis cria o cliente e faz um Ping para falhar cedo se estiver inválido.
func ConnectRedis(cfg config.RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       int(cfg.RedisDb),
		Protocol: int(cfg.RedisProtocol),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}

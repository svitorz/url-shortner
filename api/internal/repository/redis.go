package repository

import (
	"context"
	"fmt"
	"svitorz/url-shortner/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisOptions: ContextTimeoutEnabled faz o cliente respeitar o prazo do
// context (sem isso o go-redis ignora o deadline e espera o ReadTimeout de
// 3 s). O cache de links depende disso para cair no banco rápido.
func redisOptions(cfg config.RedisConfig, addr string) *redis.Options {
	return &redis.Options{
		Addr:                  addr,
		Password:              cfg.RedisPassword,
		DB:                    int(cfg.RedisDb),
		Protocol:              int(cfg.RedisProtocol),
		ContextTimeoutEnabled: true,
	}
}

// ConnectRedis cria o cliente e faz um Ping para falhar cedo se estiver inválido.
func ConnectRedis(cfg config.RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)

	rdb := redis.NewClient(redisOptions(cfg, addr))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}

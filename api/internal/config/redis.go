package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	RedisHost     string
	RedisPort     int64
	RedisPassword string
	RedisDb       int64
	RedisProtocol int64
}

func LoadRedisConf() (*RedisConfig, error) {
	var cfg *RedisConfig
	redisport, err := strconv.ParseInt(os.Getenv("REDIS_PORT"), 10, 64)

	if err != nil {
		return cfg, err
	}

	redisdb, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)

	if err != nil {
		return cfg, err
	}

	redprotocol, err := strconv.ParseInt(os.Getenv("REDIS_PROTOCOL"), 10, 64)

	if err != nil {
		return cfg, err
	}
	return &RedisConfig{
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     redisport,
		RedisDb:       redisdb,
		RedisProtocol: redprotocol,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}, nil
}

func ConnectRedis(cfg *RedisConfig) (*redis.Client, error) {

	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	fmt.Println("Debug redis ADDR:", addr)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       int(cfg.RedisDb),
		Protocol: int(cfg.RedisProtocol),
	})

	// Testa conexão (Ping) para falhar cedo se estiver inválida
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

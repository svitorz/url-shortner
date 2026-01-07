package repository

import (
	"fmt"
	"svitorz/url-shortner/internal/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func init() {
	cfg, err := config.LoadRedisConf()
	if err != nil {
		fmt.Println("Erro ao buscar configuracoes do redis")
	}

	RDB, err = config.ConnectRedis(cfg)
	if err != nil {
		fmt.Println("Erro ao se conectar redis")
	}
}

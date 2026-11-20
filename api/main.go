package main

import (
	"fmt"
	"os"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/config"
	"svitorz/url-shortner/internal/controllers"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/repository"
	"svitorz/url-shortner/internal/requests"
	"svitorz/url-shortner/internal/requests/ratelimit"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "erro fatal:", err)
		os.Exit(1)
	}
}

// run monta todas as dependências explicitamente, na ordem, e devolve o
// primeiro erro: o processo não sobe com config, banco ou Redis inválidos.
func run() error {
	fmt.Println("Iniciando execucao")

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	fmt.Println("carregou configuracoes")

	db, err := repository.ConnectDatabase(cfg.Db)
	if err != nil {
		return fmt.Errorf("conectar ao Postgres: %w", err)
	}

	rdb, err := repository.ConnectRedis(cfg.Redis)
	if err != nil {
		return fmt.Errorf("conectar ao Redis: %w", err)
	}
	defer rdb.Close()

	if err := models.LoadMigrations(db); err != nil {
		return fmt.Errorf("executar migrations: %w", err)
	}
	fmt.Println("Executando migrations")

	ctl := &controllers.Controller{
		DB:       db,
		Tokens:   auth.NewTokens(cfg.Api.ApiKey, cfg.Api.TokenLifeSpan),
		HashCost: cfg.HashCost,
		BaseURL:  cfg.PublicBaseURL,
	}

	// BUG conhecido (H10), mantido de propósito nesta fase para o refactor não
	// mudar comportamento: a janela depende do minuto do relógio e é medida em
	// nanossegundos. Corrigido na fase 2, com teste.
	window := time.Duration(cfg.RateLimit.Minutes * time.Now().Minute())
	rl := ratelimit.New(rdb, cfg.RateLimit.Limit, window)

	r := requests.SetupRouter(ctl, rl)
	fmt.Println("setup router")

	return r.Run(fmt.Sprintf(":%d", cfg.AppPort))
}

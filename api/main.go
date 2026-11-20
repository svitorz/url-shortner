package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"svitorz/url-shortner/internal/auth"
	"svitorz/url-shortner/internal/clicks"
	"svitorz/url-shortner/internal/config"
	"svitorz/url-shortner/internal/controllers"
	"svitorz/url-shortner/internal/linkcache"
	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/repository"
	"svitorz/url-shortner/internal/requests"
	"svitorz/url-shortner/internal/requests/ratelimit"
	"syscall"
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
		Cache:    linkcache.New(rdb, cfg.LinkCache.TTL(), cfg.LinkCache.NegativeTTL()),
		Clicks:   clicks.New(rdb, db),
	}

	rl := ratelimit.New(rdb, cfg.RateLimit.Limit, cfg.RateLimit.Window())

	r, err := requests.SetupRouter(ctl, rl, cfg.TrustedProxies)
	if err != nil {
		return fmt.Errorf("configurar router: %w", err)
	}
	fmt.Println("setup router")

	// SIGTERM (docker stop) e Ctrl+C encerram de forma graciosa: param de
	// aceitar conexões, terminam as requisições em andamento e fazem um último
	// flush dos cliques.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	flushDone := make(chan struct{})
	go func() {
		ctl.Clicks.Run(ctx, time.Duration(cfg.ClickFlushSeconds)*time.Second)
		close(flushDone)
	}()

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.AppPort),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.ListenAndServe() }()

	select {
	case err := <-serveErr:
		stop()
		<-flushDone
		return err
	case <-ctx.Done():
	}

	fmt.Println("encerrando")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	err = srv.Shutdown(shutdownCtx)
	<-flushDone
	return err
}

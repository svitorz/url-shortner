package main

import (
	"fmt"
	"log"
	"svitorz/url-shortner/internal/config"
	"svitorz/url-shortner/internal/requests"
	"time"
)

func main() {
	fmt.Println("Iniciando execucao")
	cfg, err := config.LoadConfig()
	fmt.Println("Carregando configuracoes")
	if err != nil {
		fmt.Println("Erro ao carregar arquivos do .env", err)
	}

	fmt.Println(cfg)

	rdb, err := config.ConnectRedis(cfg)

	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()

	r := requests.SetupRouter(rdb, 10, time.Minute)
	r.Run(fmt.Sprintf(":%d", cfg.AppPort))
}

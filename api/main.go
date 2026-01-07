package main

import (
	"fmt"
	"svitorz/url-shortner/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Iniciando execucao")
	cfg, err := config.LoadConfig()
	fmt.Println("Carregando configuracoes")
	if err != nil {
		fmt.Println("Erro ao carregar arquivos do .env", err)
	}

	fmt.Println(cfg)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		fmt.Println(c.Request.Body)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	fmt.Println("Router carregado", router)

	router.Run(fmt.Sprintf(":%d", cfg.AppPort))
}

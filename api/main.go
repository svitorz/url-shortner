package main

import (
	"fmt"
	"svitorz/url-shortner/internal/config"
	"svitorz/url-shortner/internal/requests"
)

func main() {
	fmt.Println("Iniciando execucao")
	r := requests.SetupRouter()
	fmt.Println("setup router")
	port, err := config.LoadAppConfig()
	if err != nil {
		fmt.Println("Erro ao carregar confg da aplicacao", err)
	}

	fmt.Println("carregou configuracoes")
	r.Run(fmt.Sprintf(":%d", port))
}

package repository

import (
	"fmt"
	"svitorz/url-shortner/internal/config"

	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	cfg, err := config.LoadDbConfig()
	if err != nil {
		fmt.Println("Erro ao buscar configuracoes")
	}
	DB, err = config.ConnectDatabase(cfg)
	if err != nil {
		fmt.Println("Failed to connect to database!", err)
	}
}

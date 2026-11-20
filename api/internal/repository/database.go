package repository

import (
	"fmt"
	"svitorz/url-shortner/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase abre a conexão com o Postgres. O gorm.Open faz ping por
// padrão, então um banco inacessível vira erro aqui, no boot.
func ConnectDatabase(cfg config.DbConfig) (*gorm.DB, error) {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=America/Sao_Paulo", cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dns,
		PreferSimpleProtocol: false,
	}), &gorm.Config{})
}

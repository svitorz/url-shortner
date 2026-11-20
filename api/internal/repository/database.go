package repository

import (
	"fmt"
	"log"
	"os"
	"svitorz/url-shortner/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDatabase abre a conexão com o Postgres. O gorm.Open faz ping por
// padrão, então um banco inacessível vira erro aqui, no boot.
func ConnectDatabase(cfg config.DbConfig) (*gorm.DB, error) {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=America/Sao_Paulo", cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dns,
		PreferSimpleProtocol: false,
	}), &gorm.Config{Logger: dbLogger})
}

// dbLogger loga erros e queries lentas SEM os valores dos parâmetros (o
// padrão do GORM imprimia e-mails e hashes de senha) e não trata "registro
// não encontrado" como erro, pois é fluxo normal (login, redirect).
var dbLogger = logger.New(log.New(os.Stdout, "", log.LstdFlags), logger.Config{
	SlowThreshold:             200 * time.Millisecond,
	LogLevel:                  logger.Warn,
	IgnoreRecordNotFoundError: true,
	ParameterizedQueries:      true,
	Colorful:                  false,
})

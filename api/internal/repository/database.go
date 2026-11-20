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
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dns,
		PreferSimpleProtocol: false,
	}), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return nil, err
	}
	if err := configurePool(db, cfg.MaxConns); err != nil {
		return nil, err
	}
	return db, nil
}

// configurePool deixa até n conexões abertas E ociosas. O padrão do
// database/sql é ilimitado abertas e só 2 ociosas: sob concorrência, cada
// conexão excedente é fechada ao voltar ao pool e reaberta na consulta
// seguinte (fork no Postgres + autenticação SCRAM a cada vez).
func configurePool(db *gorm.DB, n int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(n)
	sqlDB.SetMaxIdleConns(n)
	// Devolve ao Postgres as conexões que sobraram de um pico.
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	return nil
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

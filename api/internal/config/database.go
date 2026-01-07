package config

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	DbName     string
	DbHost     string
	DbPort     int64
	DbUser     string
	DbPassword string
}

func ConnectDatabase(cfg *DbConfig) (db *gorm.DB, err error) {
	dns := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=America/Sao_Paulo", cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dns,
		PreferSimpleProtocol: false,
	}), &gorm.Config{})
	return db, err
}

func LoadDbConfig() (*DbConfig, error) {
	var cfg *DbConfig

	dbport, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)

	if err != nil {
		return cfg, err
	}

	if err != nil {
		return cfg, err
	}

	return &DbConfig{
		DbName:     os.Getenv("DB_NAME"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     dbport,
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
	}, nil
}
